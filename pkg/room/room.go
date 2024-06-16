package room

import (
	"encoding/json"
	"sync"

	"github.com/Harshitk-cp/rtmp_server/pkg/errors"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
)

type Participant struct {
	ID             string
	Conn           *websocket.Conn
	PeerConnection *webrtc.PeerConnection
	mutex          sync.RWMutex
	VideoTrack     *webrtc.TrackLocalStaticRTP
	AudioTrack     *webrtc.TrackLocalStaticRTP
}

type Room struct {
	ID              string
	Participants    map[string]*Participant
	ServerPeer      *Participant
	ServerPeerMutex sync.RWMutex
	mutex           sync.RWMutex
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
		ID:         id,
		Conn:       conn,
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

func (r *Room) SendToParticipant(participantID, msgType, data, fromClientID string) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	participant, exists := r.Participants[participantID]
	if !exists {
		logrus.Errorf("Participant %s not found in room %s", participantID, r.ID)
		return
	}

	message := SignalMessage{
		Type:         msgType,
		Data:         data,
		FromClientID: fromClientID,
	}
	participant.mutex.Lock()
	err := participant.Conn.WriteJSON(message)
	participant.mutex.Unlock()
	if err != nil {
		logrus.Errorf("Error sending message to participant %s: %v", participantID, err)
	}
}

func (r *Room) CreateAndBroadcastOffer() error {
	r.ServerPeerMutex.Lock()
	defer r.ServerPeerMutex.Unlock()

	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		logrus.Errorf("Error creating PeerConnection: %v", err)
		return err
	}

	videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264}, "video", "pion")
	if err != nil {
		return err
	}
	audioTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus}, "audio", "pion")
	if err != nil {
		return err
	}

	_, err = pc.AddTrack(videoTrack)
	if err != nil {
		return err
	}
	_, err = pc.AddTrack(audioTrack)
	if err != nil {
		return err
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

	r.ServerPeer = &Participant{
		ID:             "server",
		PeerConnection: pc,
		VideoTrack:     videoTrack,
		AudioTrack:     audioTrack,
	}

	pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}
		candidateJSON, err := json.Marshal(candidate.ToJSON())
		if err != nil {
			logrus.Errorf("Error marshaling ICE candidate: %v", err)
			return
		}
		r.Broadcast("server", "getCandidate", string(candidateJSON), "server")
	})

	r.Broadcast("server", "getOffer", offer.SDP, "server")
	return nil
}

func (r *Room) HandleAnswer(fromClientID, answer string) error {
	r.ServerPeerMutex.RLock()
	defer r.ServerPeerMutex.RUnlock()

	if r.ServerPeer == nil || r.ServerPeer.PeerConnection == nil {
		logrus.Error("Server peer connection not found")
		return errors.ErrPeerConnectionNotFound
	}

	answerSDP := webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  answer,
	}
	err := r.ServerPeer.PeerConnection.SetRemoteDescription(answerSDP)
	if err != nil {
		logrus.Errorf("Error setting remote description: %v", err)
		return err
	}
	return nil
}
