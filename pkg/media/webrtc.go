package media

import (
	"fmt"

	"github.com/Harshitk-cp/rtmp_server/pkg/room"
	"github.com/pion/webrtc/v3"
)

type WebRTCSignaling struct {
	participant *room.Participant
	peerConn    *webrtc.PeerConnection
}

func NewWebRTCSignaling() *WebRTCSignaling {
	return &WebRTCSignaling{}
}

func (s *WebRTCSignaling) Start(participant *room.Participant) error {
	s.participant = participant

	config := webrtc.Configuration{}
	peerConn, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return fmt.Errorf("peer connection failed: %w", err)
	}
	s.peerConn = peerConn

	s.setupPeerConnectionHandlers()
	s.negotiateWebRTCConnection()

	return nil
}

func (s *WebRTCSignaling) setupPeerConnectionHandlers() {
	s.peerConn.OnICECandidate(func(candidate *webrtc.ICECandidate) {
	})

	s.peerConn.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
	})

}

func (s *WebRTCSignaling) negotiateWebRTCConnection() {
	offer, err := s.peerConn.CreateOffer(nil)
	if err != nil {
		return
	}

	err = s.peerConn.SetLocalDescription(offer)
	if err != nil {
		return
	}

}

func (s *WebRTCSignaling) Stop() {
	s.peerConn.Close()
}
