package rtmp

import (
	"fmt"
	"net"

	"github.com/Harshitk-cp/rtmp_server/pkg/room"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media/samplebuilder"
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

	// for {
	// 	select {
	// 	case videoBuf := <-s.RTMPHandler.videoRTPChan:
	// 		err := rm.StreamingTracks.VideoTrack.WriteSample(media.Sample{
	// 			Data:     videoBuf,
	// 			Duration: 128 * time.Millisecond,
	// 		})
	// 		if err != nil {
	// 			logrus.Errorf("Error sending RTP packet to new video track: %v", err)

	// 		}
	// 	case audioBuf := <-s.RTMPHandler.audioRTPChan:
	// 		err := rm.StreamingTracks.AudioTrack.WriteSample(media.Sample{
	// 			Data:     audioBuf,
	// 			Duration: time.Second / 30,
	// 		})
	// 		if err != nil {
	// 			logrus.Errorf("Error sending RTP packet to new video track: %v", err)

	// 		}
	// 	}
	// }

	return nil
}

func (s *SFUServer) RtpToTrack(track *webrtc.TrackLocalStaticSample, depacketizer rtp.Depacketizer, sampleRate uint32, port int) {
	// Open a UDP Listener for RTP Packets on port 5004
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: port})
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = listener.Close(); err != nil {
			panic(err)
		}
	}()

	sampleBuffer := samplebuilder.New(10, depacketizer, sampleRate)

	for {
		inboundRTPPacket := make([]byte, 1500)
		packet := &rtp.Packet{}

		n, _, err := listener.ReadFrom(inboundRTPPacket)
		if err != nil {
			panic(fmt.Sprintf("error during read: %s", err))
		}

		if err = packet.Unmarshal(inboundRTPPacket[:n]); err != nil {
			panic(err)
		}

		sampleBuffer.Push(packet)
		for {
			sample := sampleBuffer.Pop()
			if sample == nil {
				break
			}

			if writeErr := track.WriteSample(*sample); writeErr != nil {
				panic(writeErr)
			}
		}
	}
}
