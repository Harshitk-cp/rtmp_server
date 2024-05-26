package service

import (
	"fmt"
	"io"
	"net"

	"github.com/Harshitk-cp/rtmp_server/pkg/config"
	"github.com/livekit/protocol/logger"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/yutopp/go-rtmp"
	rtmpmsg "github.com/yutopp/go-rtmp/message"
)

type Relay struct {
	rtmpServer *rtmp.Server
}

func NewRelay() *Relay {
	return &Relay{}
}

func (r *Relay) Start(conf *config.Config, port int) error {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Errorw("failed to start TCP listener", err, "port", port)
		return err
	}

	h := NewRTMPHandler()
	srv := rtmp.NewServer(&rtmp.ServerConfig{
		OnConnect: func(conn net.Conn) (io.ReadWriteCloser, *rtmp.ConnConfig) {
			l := log.StandardLogger()
			if conf.Logging.JSON {
				l.SetFormatter(&log.JSONFormatter{})
			}
			lf := l.WithFields(conf.GetLoggerFields())
			logrus.WithFields(logrus.Fields{
				"remoteAddr": conn.RemoteAddr(),
			}).Info("Client connected")
			return conn, &rtmp.ConnConfig{
				Handler: h,
				ControlState: rtmp.StreamControlStateConfig{
					DefaultBandwidthWindowSize: 6 * 1024 * 1024 / 8,
				},
				Logger: lf,
			}
		},
	})

	r.rtmpServer = srv
	log.Printf("RTMP relay server started on %s", addr)
	go func() {
		if err := r.rtmpServer.Serve(listener); err != nil {
			log.Printf("RTMP relay server stopped: %v", err)
		}
	}()

	return nil
}

func (s *Relay) Stop() error {
	return s.rtmpServer.Close()
}

type RTMPHandler struct {
	rtmp.DefaultHandler

	log logger.Logger
}

func NewRTMPHandler() *RTMPHandler {
	h := &RTMPHandler{
		log: logger.GetLogger(),
	}
	return h
}

func (h *RTMPHandler) OnPublish(_ *rtmp.StreamContext, timestamp uint32, cmd *rtmpmsg.NetStreamPublish) error {
	logrus.Infof("OnPublish called with publishing name: %s", cmd.PublishingName)
	if cmd.PublishingName != "streamix" {
		err := fmt.Errorf("unsupported publishing name: %s", cmd.PublishingName)
		logrus.Error(err)
		return err
	}
	logrus.Info("Publishing name is supported")
	return nil
}

// func (h *RTMPHandler) OnVideo(timestamp uint32, payload io.Reader) error {
// 	logrus.Infof("Received video data: timestamp=%d", timestamp)

// 	var video flvtag.VideoData
// 	if err := flvtag.DecodeVideoData(payload, &video); err != nil {
// 		return err
// 	}

// 	flvBody := new(bytes.Buffer)
// 	if _, err := io.Copy(flvBody, video.Data); err != nil {
// 		return err
// 	}

// 	log.Printf("FLV Video Data: Timestamp = %d, FrameType = %+v, CodecID = %+v, AVCPacketType = %+v, CT = %+v, Data length = %+v",
// 		timestamp,
// 		video.FrameType,
// 		video.CodecID,
// 		video.AVCPacketType,
// 		video.CompositionTime,
// 		len(flvBody.Bytes()),
// 	)

// 	// You can process the video data here if needed
// 	return nil
// }

// func (h *RTMPHandler) OnAudio(timestamp uint32, payload io.Reader) error {
// 	logrus.Infof("Received audio data: timestamp=%d", timestamp)

// 	var audio flvtag.AudioData
// 	if err := flvtag.DecodeAudioData(payload, &audio); err != nil {
// 		return err
// 	}

// 	flvBody := new(bytes.Buffer)
// 	if _, err := io.Copy(flvBody, audio.Data); err != nil {
// 		return err
// 	}

// 	log.Printf("FLV Audio Data: Timestamp = %d, SoundFormat = %+v, SoundRate = %+v, SoundSize = %+v, SoundType = %+v, AACPacketType = %+v, Data length = %+v",
// 		timestamp,
// 		audio.SoundFormat,
// 		audio.SoundRate,
// 		audio.SoundSize,
// 		audio.SoundType,
// 		audio.AACPacketType,
// 		len(flvBody.Bytes()),
// 	)

// 	// You can process the audio data here if needed
// 	return nil
// }
