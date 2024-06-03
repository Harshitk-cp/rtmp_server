package media

import (
	"fmt"

	"github.com/Harshitk-cp/rtmp_server/pkg/room"
	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
)

type WebRTCSignaling struct {
	participant *room.Participant
	peerConn    *webrtc.PeerConnection
}

func NewWebRTCSignaling() *WebRTCSignaling {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
		SDPSemantics: webrtc.SDPSemanticsUnifiedPlan,
	}

	mediaEngine := webrtc.MediaEngine{}
	if err := mediaEngine.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:  webrtc.MimeTypeOpus,
			ClockRate: 48000,
			Channels:  2,
		},
		PayloadType: 111,
	}, webrtc.RTPCodecTypeAudio); err != nil {
		logrus.Errorf("failed to register Opus codec: %v", err)
		return nil
	}

	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))
	peerConn, err := api.NewPeerConnection(config)
	if err != nil {
		logrus.Errorf("peer connection failed: %v", err)
		return nil
	}

	return &WebRTCSignaling{
		peerConn: peerConn,
	}
}

func (s *WebRTCSignaling) Start(participant *room.Participant, offerSDP string) (string, error) {
	s.participant = participant

	videoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "video/vp8"}, "video", "pion")
	if err != nil {
		return "", fmt.Errorf("failed to create video track: %w", err)
	}
	rtpSender, err := s.peerConn.AddTrack(videoTrack)
	if err != nil {
		return "", fmt.Errorf("failed to add video track: %w", err)
	}
	processRTCP(rtpSender)

	audioTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "audio/opus"}, "audio", "pion")
	if err != nil {
		return "", fmt.Errorf("failed to create audio track: %w", err)
	}
	rtpSender, err = s.peerConn.AddTrack(audioTrack)
	if err != nil {
		return "", fmt.Errorf("failed to add audio track: %w", err)
	}
	processRTCP(rtpSender)

	s.setupPeerConnectionHandlers()

	answerSDP, err := s.negotiateWebRTCConnection(string(offerSDP))
	if err != nil {
		return "", fmt.Errorf("failed to negotiate WebRTC connection: %w", err)
	}

	logrus.Infof("Answer SDP: %s", answerSDP)
	// go rtpToTrack(videoTrack, &codecs.VP8Packet{}, 90000, 5004)
	// go rtpToTrack(audioTrack, &codecs.OpusPacket{}, 48000, 5006)

	return answerSDP, nil
}

func (s *WebRTCSignaling) setupPeerConnectionHandlers() {
	s.peerConn.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	})

	s.peerConn.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			logrus.Info("ICE candidate gathering complete")
		}
	})

	s.peerConn.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		logrus.Infof("Received track: %s", track.Kind())
	})
}

func (s *WebRTCSignaling) negotiateWebRTCConnection(offerSDP string) (string, error) {
	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  offerSDP,
	}

	if err := s.peerConn.SetRemoteDescription(offer); err != nil {
		return "", fmt.Errorf("failed to set remote description: %w", err)
	}

	answer, err := s.peerConn.CreateAnswer(nil)
	if err != nil {
		return "", fmt.Errorf("failed to create answer: %w", err)
	}

	gatherComplete := webrtc.GatheringCompletePromise(s.peerConn)

	if err := s.peerConn.SetLocalDescription(answer); err != nil {
		return "", fmt.Errorf("failed to set local description: %w", err)
	}

	<-gatherComplete

	return s.peerConn.LocalDescription().SDP, nil
}

func (s *WebRTCSignaling) Stop() {
	s.peerConn.Close()
}

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
