package rtmp

import (
	"encoding/json"

	"github.com/Harshitk-cp/rtmp_server/pkg/errors"
	"github.com/Harshitk-cp/rtmp_server/pkg/room"

	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
)

type SFUServer struct {
	RTMPHandler *RTMPHandler
	RoomManager *room.RoomManager
}

func NewSFUServer(rm *room.RoomManager, rtmpHandler *RTMPHandler) *SFUServer {
	return &SFUServer{
		RTMPHandler: rtmpHandler,
		RoomManager: rm,
	}
}

func (s *SFUServer) HandleOffer(participant *room.Participant, offer string, rm *room.Room) string {

	// Create a new PeerConnection
	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		logrus.Errorf("Error creating PeerConnection: %v", err)
		return ""
	}

	// Set the remote description
	err = pc.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  offer,
	})
	if err != nil {
		logrus.Errorf("Error setting remote description: %v", err)
		return ""
	}

	// Add video and audio tracks to the PeerConnection
	if participant.VideoTrack != nil {
		_, err := pc.AddTrack(participant.VideoTrack)
		if err != nil {
			logrus.Errorf("Error adding video track: %v", err)
			return ""
		}
	}
	if participant.AudioTrack != nil {
		_, err := pc.AddTrack(participant.AudioTrack)
		if err != nil {
			logrus.Errorf("Error adding audio track: %v", err)
			return ""
		}
	}

	// Create answer
	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		logrus.Errorf("Error creating answer: %v", err)
		return ""
	}

	// Set the local description
	err = pc.SetLocalDescription(answer)
	if err != nil {
		logrus.Errorf("Error setting local description: %v", err)
		return ""
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
		rm.SendToParticipant(participant.ID, "getCandidate", string(candidateJSON), participant.ID)
	})

	// Set the PeerConnection on the participant
	participant.PeerConnection = pc

	// Send the answer back to the client
	return answer.SDP
}

func (s *SFUServer) HandleCandidate(participant *room.Participant, candidate string) {
	if participant.PeerConnection == nil {
		logrus.Error("PeerConnection not found for participant")
		return
	}

	err := participant.PeerConnection.AddICECandidate(webrtc.ICECandidateInit{
		Candidate: candidate,
		// SDPMid:        *sdpMid,
		// SDPMLineIndex: sdpMLineIndex,
	})
	if err != nil {
		logrus.Errorf("Error adding ICE candidate: %v", err)
	}
}

func (s *SFUServer) SendRTMPToWebRTC(participant *room.Participant) error {
	const maxPacketSize = 150000
	for {
		select {
		case rtpPacket := <-s.RTMPHandler.videoRTPChan:
			if len(rtpPacket.Payload) <= maxPacketSize {
				err := participant.VideoTrack.WriteRTP(rtpPacket)
				if err != nil {
					logrus.Errorf("Error sending RTP packet to video track: %v", err)
					return errors.ErrConvertRtmpToWebrtc
				}
			} else {
				logrus.Warnf("Skipping large video RTP packet of size: %d", len(rtpPacket.Payload))
			}
		case rtpPacket := <-s.RTMPHandler.audioRTPChan:
			err := participant.AudioTrack.WriteRTP(rtpPacket)
			if err != nil {
				logrus.Errorf("Error sending RTP packet to audio track: %v", err)
				return errors.ErrConvertRtmpToWebrtc
			}
		}
	}
}
