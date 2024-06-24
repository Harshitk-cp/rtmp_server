package rtmp

import (
	"sync"
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

func (s *SFUServer) SendRTMPToWebRTC(rm *room.Room) error {
	var (
		ingressActive bool
		mu            sync.Mutex
		lastActivity  time.Time
	)

	inactivityThreshold := 5 * time.Second

	updateActivity := func() {
		mu.Lock()
		lastActivity = time.Now()
		if !ingressActive {
			ingressActive = true
			s.RTMPHandler.webhookManager.SendWebhook("ingress_started", rm.IngressID)
			logrus.Infof("Ingress started for room %s", rm.IngressID)
			rm.IsPublishing = true
			rm.Broadcast("server", "isPublishing", string("true"), rm.ID)
		}
		mu.Unlock()
	}

	checkInactivity := func() {
		mu.Lock()
		if ingressActive && time.Since(lastActivity) > inactivityThreshold {
			ingressActive = false
			s.RTMPHandler.webhookManager.SendWebhook("ingress_ended", rm.IngressID)
			rm.IsPublishing = false
			rm.Broadcast("server", "isPublishing", string("false"), rm.ID)
		}
		mu.Unlock()
	}

	go func() {
		for {
			select {
			case videoBuf, ok := <-s.RTMPHandler.videoRTPChan:
				if !ok {
					return
				}
				updateActivity()
				err := rm.StreamingTracks.VideoTrack.WriteSample(media.Sample{
					Data:     videoBuf,
					Duration: 128 * time.Millisecond,
				})

				if err != nil {
					logrus.Errorf("Error sending RTP packet to video track: %v", err)
				}
			case audioBuf, ok := <-s.RTMPHandler.audioRTPChan:
				if !ok {
					return
				}
				updateActivity()
				err := rm.StreamingTracks.AudioTrack.WriteSample(media.Sample{
					Data:     audioBuf,
					Duration: time.Second / 60,
				})
				if err != nil {
					logrus.Errorf("Error sending RTP packet to audio track: %v", err)
				}
			case <-time.After(time.Second):
				checkInactivity()
			}
		}
	}()

	return nil
}
