package rtmp

import (
	"github.com/Harshitk-cp/rtmp_server/pkg/room"
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

func (s *SFUServer) SendRTMPToWebRTC(participant *room.Participant) {
	for {
		select {
		case rtpPacket := <-s.RTMPHandler.videoRTPChan:
			err := participant.VideoTrack.WriteRTP(rtpPacket)
			if err != nil {
				logrus.Errorf("Error sending RTP packet to video track: %v", err)
				return
			}
		case rtpPacket := <-s.RTMPHandler.audioRTPChan:
			err := participant.AudioTrack.WriteRTP(rtpPacket)
			if err != nil {
				logrus.Errorf("Error sending RTP packet to audio track: %v", err)
				return
			}
		}
	}
}
