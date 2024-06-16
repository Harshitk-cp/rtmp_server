package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Harshitk-cp/rtmp_server/pkg/room"
	"github.com/Harshitk-cp/rtmp_server/pkg/rtmp"

	"github.com/gorilla/websocket"
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
	ClientID     string `json:"clientID"`
	FromClientID string `json:"fromClientID"`
	RoomID       string `json:"roomID"`
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

		participant, err := rm.CreateParticipant(clientID, conn)
		if err != nil {
			logrus.Errorf("Error creating participant: %v", err)
			http.Error(w, "Failed to create participant", http.StatusInternalServerError)
			return
		}

		go sfuServer.SendRTMPToWebRTC(participant)

		defer func() {
			rm.RemoveParticipant(clientID)
			logrus.Infof("Participant %s removed from room %s", clientID, roomID)
		}()

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
			case "offer":
				handleOffer(rm, participant, signal.Data, signal.FromClientID)
			case "answer":
				handleAnswer(rm, signal.FromClientID, signal.Data)
			case "candidate":
				handleCandidate(rm, signal.Data, signal.FromClientID)
			default:
				logrus.Warnf("Unknown message type: %v", signal.Type)
			}
		}
	}
}

func handleOffer(rm *room.Room, participant *room.Participant, offer, fromClientId string) {
	rm.Broadcast(participant.ID, "getOffer", offer, fromClientId)
}

func handleAnswer(rm *room.Room, fromClientID, answer string) {
	rm.SendToParticipant(fromClientID, "getAnswer", answer)
}

func handleCandidate(rm *room.Room, candidate, fromClientId string) {
	rm.Broadcast("", "getCandidate", candidate, fromClientId)
}
