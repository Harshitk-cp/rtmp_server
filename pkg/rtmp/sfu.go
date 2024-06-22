package rtmp

import (
	"time"

	"github.com/Harshitk-cp/rtmp_server/pkg/room"
	"github.com/pion/webrtc/v3/pkg/media"
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

func (s *SFUServer) SendRTMPToWebRTC(rm *room.Room, clientID string) error {
	go func() {
		for videoBuf := range s.RTMPHandler.videoRTPChan {
			err := rm.StreamingTracks.VideoTrack.WriteSample(media.Sample{
				Data:     videoBuf,
				Duration: 128 * time.Millisecond,
			})
			if err != nil {
				logrus.Errorf("Error sending RTP packet to video track: %v", err)
			}
		}
	}()

	go func() {
		for audioBuf := range s.RTMPHandler.audioRTPChan {
			err := rm.StreamingTracks.AudioTrack.WriteSample(media.Sample{
				Data:     audioBuf,
				Duration: time.Second / 60,
			})
			if err != nil {
				logrus.Errorf("Error sending RTP packet to audio track: %v", err)
			}
		}
	}()

	return nil
}
