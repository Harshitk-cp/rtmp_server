package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Harshitk-cp/rtmp_server/pkg/room"
	"github.com/Harshitk-cp/rtmp_server/pkg/rtmp"
	"github.com/gorilla/websocket"
	"github.com/pion/rtp/codecs"
	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type SignalMessage struct {
	Type         string `json:"type"`
	Data         string `json:"data"`
	FromClientID string `json:"fromClientID"`
}

func WebSocketHandler(sfuServer *rtmp.SFUServer, roomManager *room.RoomManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to set websocket upgrade: %v", err)
			return
		}
		defer conn.Close()

		clientID := r.URL.Query().Get("clientID")
		roomID := r.URL.Query().Get("roomID")

		if clientID == "" || roomID == "" {
			http.Error(w, "Missing required parameters", http.StatusBadRequest)
			return
		}

		rm, exist := roomManager.GetRoom(roomID)
		if !exist {
			rm = roomManager.CreateRoom(roomID)
		}

		pc, err := webrtc.NewPeerConnection(webrtc.Configuration{
			ICEServers: []webrtc.ICEServer{
				{
					URLs: []string{"stun:stun.l.google.com:19302"},
				},
			},
		})
		if err != nil {
			logrus.Errorf("Error creating PeerConnection: %v", err)
			http.Error(w, "Failed to create PeerConnection", http.StatusInternalServerError)
			return
		}

		rtpSender, err := pc.AddTrack(rm.StreamingTracks.VideoTrack)
		if err != nil {
			logrus.Errorf("Error adding video track: %v", err)
			http.Error(w, "Failed to add video track", http.StatusInternalServerError)
			return
		}

		processRTCP(rtpSender)

		rtpSender, err = pc.AddTrack(rm.StreamingTracks.AudioTrack)
		if err != nil {
			logrus.Errorf("Error adding audio track: %v", err)
			http.Error(w, "Failed to add audio track", http.StatusInternalServerError)
			return
		}

		processRTCP(rtpSender)

		pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
			if candidate == nil {
				return
			}
			candidateJSON, err := json.Marshal(candidate.ToJSON())
			if err != nil {
				logrus.Errorf("Error marshaling ICE candidate: %v", err)
				return
			}
			rm.Broadcast("server", "getCandidate", string(candidateJSON), "server")
			message := SignalMessage{
				Type:         "candidate",
				Data:         string(candidateJSON),
				FromClientID: "server",
			}
			logrus.Info(message)

		})

		_, err = rm.CreateParticipant(clientID, conn)
		if err != nil {
			logrus.Errorf("Error creating participant: %v", err)
			http.Error(w, "Failed to create participant", http.StatusInternalServerError)
			return
		}

		rm.ClientPeerConnections[clientID] = pc

		defer func() {
			rm.RemoveParticipant(clientID)
			logrus.Infof("Participant %s removed from room %s", clientID, roomID)
		}()

		err = rm.GenerateOffer(clientID)
		if err != nil {
			logrus.Errorf("Error generating offer: %v", err)
		}

		gatherComplete := webrtc.GatheringCompletePromise(pc)

		<-gatherComplete

		// go sfuServer.SendRTMPToWebRTC(rm, clientID)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				rm.RemoveParticipant(clientID)
				break
			}

			var signal SignalMessage
			err = json.Unmarshal(msg, &signal)
			if err != nil {
				logrus.Errorf("Invalid message format: %v", err)
				continue
			}

			switch signal.Type {
			case "answer":
				handleClientAnswer(rm, signal.Data, clientID)
			case "candidate":
				handleCandidate(rm, signal.Data, clientID)
			case "startStream":
				handleStartStream(sfuServer, rm)
			default:
				logrus.Warnf("Unknown message type: %v", signal.Type)
			}
		}
	}
}

func handleClientAnswer(rm *room.Room, answer, clientID string) {
	rm.HandleAnswer(answer, clientID)
}

func handleCandidate(rm *room.Room, candidate, clientID string) {
	pc, ok := rm.ClientPeerConnections[clientID]
	if !ok {
		logrus.Errorf("Error peer connection not found for client ID: %v", clientID)
	}

	err := pc.AddICECandidate(webrtc.ICECandidateInit{
		Candidate: candidate,
	})
	logrus.Infof("Added candidate: %v", candidate)
	if err != nil {
		logrus.Errorf("Error adding ICE candidate: %v", err)
	}

}

func handleStartStream(sfuServer *rtmp.SFUServer, rm *room.Room) {
	log.Printf("Starting RTP to Track with %v users", len(rm.Participants))
	go sfuServer.RtpToTrack(rm.StreamingTracks.VideoTrack, &codecs.VP8Packet{}, 90000, 5007)
	go sfuServer.RtpToTrack(rm.StreamingTracks.AudioTrack, &codecs.OpusPacket{}, 48000, 5009)
}

// like NACK this needs to be called.
func processRTCP(rtpSender *webrtc.RTPSender) {
	go func() {
		rtcpBuf := make([]byte, 1500)

		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()
}
