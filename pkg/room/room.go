package room

import (
	"sync"

	"github.com/Harshitk-cp/rtmp_server/pkg/errors"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
)

type Participant struct {
	ID         string
	Conn       *websocket.Conn
	PeerConn   *webrtc.PeerConnection
	VideoTrack *webrtc.TrackLocalStaticRTP
	AudioTrack *webrtc.TrackLocalStaticRTP
	mutex      sync.RWMutex
}

type Room struct {
	ID           string
	Participants map[string]*Participant
	mutex        sync.RWMutex
}

func NewRoom(id string) *Room {
	return &Room{
		ID:           id,
		Participants: make(map[string]*Participant),
	}
}

func (r *Room) CreateParticipant(id string, conn *websocket.Conn) (*Participant, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.Participants[id]; exists {
		return nil, errors.ErrParticipantExists
	}

	videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264}, "video", "pion")
	if err != nil {
		return nil, err
	}

	audioTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus}, "audio", "pion")
	if err != nil {
		return nil, err
	}

	participant := &Participant{
		ID:   id,
		Conn: conn,
		// PeerConn:   peerConn,
		VideoTrack: videoTrack,
		AudioTrack: audioTrack,
	}

	r.Participants[id] = participant
	return participant, nil
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
	if participant.PeerConn != nil {
		participant.PeerConn.Close()
	}
	return nil
}

type SignalMessage struct {
	Type         string `json:"type"`
	Data         string `json:"data"`
	FromClientID string `json:"fromClientID"`
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
			participant.mutex.Lock()
			err := participant.Conn.WriteJSON(message)
			participant.mutex.Unlock()
			if err != nil {
				logrus.Errorf("Error broadcasting message to participant %s: %v", participant.ID, err)
			}
		}
	}
}

func (r *Room) SendToParticipant(participantID, msgType, data string) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	participant, exists := r.Participants[participantID]
	if !exists {
		logrus.Errorf("Participant %s not found in room %s", participantID, r.ID)
		return
	}

	message := SignalMessage{
		Type: msgType,
		Data: data,
	}
	participant.mutex.Lock()
	err := participant.Conn.WriteJSON(message)
	participant.mutex.Unlock()
	if err != nil {
		logrus.Errorf("Error sending message to participant %s: %v", participantID, err)
	}
}
