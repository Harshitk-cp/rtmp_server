package room

import (
	"encoding/json"
	"io"
	"sync"

	"github.com/Harshitk-cp/rtmp_server/pkg/errors"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
)

type Participant struct {
	ID   string
	Conn *websocket.Conn

	mutex sync.RWMutex
}

func NewParticipant(id string, conn *websocket.Conn) *Participant {
	return &Participant{
		ID:   id,
		Conn: conn,
	}
}

func (p *Participant) Send(message SignalMessage) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	err := p.Conn.WriteJSON(message)
	if err != nil {
		logrus.Errorf("Error sending message to participant %v: %v", p.ID, err)
	}
}

type Room struct {
	ID                    string
	Participants          map[string]*Participant
	ClientPeerConnections map[string]*webrtc.PeerConnection
	StreamingTracks       *StreamingTracks
	mutex                 sync.RWMutex
}

func NewRoom(id string) *Room {
	streamingTracks, err := NewStreamingTracks()
	if err != nil {
		logrus.Fatalf("Error creating streaming tracks: %v", err)
	}

	return &Room{
		ID:                    id,
		Participants:          make(map[string]*Participant),
		ClientPeerConnections: make(map[string]*webrtc.PeerConnection),
		StreamingTracks:       streamingTracks,
	}
}

type StreamingTracks struct {
	VideoTrack *webrtc.TrackLocalStaticSample
	AudioTrack *webrtc.TrackLocalStaticSample
}

func NewStreamingTracks() (*StreamingTracks, error) {
	videoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264, ClockRate: 90000}, "video", "pion")
	if err != nil {
		logrus.Errorf("Error creating video track: %v", err)

	}

	audioTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus, ClockRate: 48000}, "audio", "pion")
	if err != nil {
		logrus.Errorf("Error creating audio track: %v", err)

	}

	return &StreamingTracks{
		VideoTrack: videoTrack,
		AudioTrack: audioTrack,
	}, nil
}

func (r *Room) CreateParticipant(id string, conn *websocket.Conn) (*Participant, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	p := NewParticipant(id, conn)
	r.Participants[id] = p
	return p, nil
}

func (r *Room) RemoveParticipant(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	participant, exists := r.Participants[id]
	if !exists {
		return errors.ErrParticipantNotFound
	}

	delete(r.Participants, id)
	participant.Conn.Close()
	return nil
}

type SignalMessage struct {
	Type         string `json:"type"`
	Data         string `json:"data"`
	FromClientID string `json:"fromClientID"`
}

func (r *Room) SendToParticipant(participantID, msgType, data, fromID string) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if participant, exists := r.Participants[participantID]; exists {
		message := SignalMessage{
			Type:         msgType,
			Data:         data,
			FromClientID: fromID,
		}
		participant.Send(message)
	}
}

func (rm *Room) HandleAnswer(answer, clientID string) {
	pc, ok := rm.ClientPeerConnections[clientID]
	if !ok {
		logrus.Errorf("Error peer connection not found for client ID: %v", clientID)
	}

	err := pc.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  answer,
	})
	if err != nil {
		logrus.Errorf("Error setting remote description: %v", err)
	}
}

func (rm *Room) GenerateOffer(clientID string) error {
	pc, ok := rm.ClientPeerConnections[clientID]
	if !ok {
		return errors.ErrPeerConnectionNotFound
	}

	offer, err := pc.CreateOffer(nil)
	if err != nil {
		logrus.Errorf("Error creating offer: %v", err)
		return err
	}

	err = pc.SetLocalDescription(offer)
	if err != nil {
		logrus.Errorf("Error setting local description: %v", err)
		return err
	}

	participant, ok := rm.Participants[clientID]
	if !ok {
		return errors.ErrParticipantNotFound
	}

	pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		buffer := make([]byte, 1500) // Adjust buffer size as needed
		for {
			n, _, err := track.Read(buffer)
			if err != nil {
				if err == io.EOF {
					logrus.Infof("Track %s reached EOF", track.ID())
					break
				}
				logrus.Errorf("Error reading track %s: %v", track.ID(), err)
				continue
			}
			logrus.Infof("Track %s data: %v", track.ID(), buffer[:n])
		}
	})

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
	})

	message := SignalMessage{
		Type:         "getOffer",
		Data:         offer.SDP,
		FromClientID: "server",
	}
	participant.Send(message)

	return nil
}

func (r *Room) Broadcast(senderID, msgType, data, fromClientID string) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for id, participant := range r.Participants {
		if id != senderID {
			message := SignalMessage{
				Type:         msgType,
				Data:         data,
				FromClientID: fromClientID,
			}
			participant.Send(message)

		}
	}
}
