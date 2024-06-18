package rtmp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"path"

	"github.com/livekit/protocol/logger"
	protoutils "github.com/livekit/protocol/utils"
	log "github.com/sirupsen/logrus"
	flvtag "github.com/yutopp/go-flv/tag"
	"github.com/yutopp/go-rtmp"
	rtmpmsg "github.com/yutopp/go-rtmp/message"

	"github.com/Harshitk-cp/rtmp_server/pkg/config"
	"github.com/Harshitk-cp/rtmp_server/pkg/errors"
	"github.com/Harshitk-cp/rtmp_server/pkg/params"
)

type RTMPServer struct {
	server   *rtmp.Server
	sfu      *SFUServer
	handlers map[string]*RTMPHandler
}

func NewRTMPServer(sfu *SFUServer) *RTMPServer {
	return &RTMPServer{
		sfu:      sfu,
		handlers: make(map[string]*RTMPHandler),
	}
}

func (s *RTMPServer) Start(conf *config.Config, h *RTMPHandler, onPublish func(streamKey, resourceId string) (*params.Params, error)) error {
	port := conf.RTMPPort

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		logger.Errorw("failed to start TCP listener", err, "port", port)
		return err
	}

	srv := rtmp.NewServer(&rtmp.ServerConfig{
		OnConnect: func(conn net.Conn) (io.ReadWriteCloser, *rtmp.ConnConfig) {
			l := log.StandardLogger()

			s.handlers[h.resourceId] = h

			return conn, &rtmp.ConnConfig{
				Handler: h,
				Logger:  l,
			}
		},
	})

	s.server = srv
	go func() {
		if err := srv.Serve(listener); err != nil {
			logger.Errorw("failed to start RTMP server", err)
		}
	}()

	return nil
}

func (s *RTMPServer) CloseHandler(resourceId string) {
	h, ok := s.handlers[resourceId]
	if ok && h != nil {
		delete(s.handlers, resourceId)
	}
}

func (s *RTMPServer) Stop() error {
	return s.server.Close()
}

type RTMPHandler struct {
	rtmp.DefaultHandler

	params     *params.Params
	resourceId string
	log        logger.Logger

	onPublish func(streamKey, resourceId string) (*params.Params, error)
	onClose   func(resourceId string)

	videoRTPChan chan []byte
	audioRTPChan chan []byte
}

func NewRTMPHandler() *RTMPHandler {
	return &RTMPHandler{
		videoRTPChan: make(chan []byte, 100),
		audioRTPChan: make(chan []byte, 100),
	}
}

func (h *RTMPHandler) OnPublishCallback(cb func(streamKey, resourceId string) (*params.Params, error)) {
	h.onPublish = cb
}

func (h *RTMPHandler) OnCloseCallback(cb func(resourceId string)) {
	h.onClose = cb
}

func (h *RTMPHandler) OnPublish(_ *rtmp.StreamContext, timestamp uint32, cmd *rtmpmsg.NetStreamPublish) error {
	if cmd.PublishingName == "" {
		return errors.ErrMissingStreamKey
	}

	_, streamKey := path.Split(cmd.PublishingName)
	appName := "live"
	h.resourceId = protoutils.NewGuid(protoutils.RTMPResourcePrefix)
	h.log = logger.GetLogger().WithValues("appName", appName, "streamKey", streamKey, "resourceID", h.resourceId)
	if h.onPublish != nil {
		params, err := h.onPublish(streamKey, h.resourceId)
		if err != nil {
			return err
		}
		h.params = params
	}

	log.WithFields(log.Fields{
		"appName":        appName,
		"streamKey":      streamKey,
		"resourceID":     h.resourceId,
		"publishingName": cmd.PublishingName,
	}).Info("Received publish request")

	h.log.Infow("Received a new published stream")

	return nil
}

func (h *RTMPHandler) OnAudio(timestamp uint32, payload io.Reader) error {
	var audio flvtag.AudioData
	if err := flvtag.DecodeAudioData(payload, &audio); err != nil {
		return err
	}

	data := new(bytes.Buffer)
	if _, err := io.Copy(data, audio.Data); err != nil {
		return err
	}

	h.audioRTPChan <- data.Bytes()
	return nil
}

const headerLengthField = 4

func (h *RTMPHandler) OnVideo(timestamp uint32, payload io.Reader) error {
	var video flvtag.VideoData
	if err := flvtag.DecodeVideoData(payload, &video); err != nil {
		return err
	}

	data := new(bytes.Buffer)
	if _, err := io.Copy(data, video.Data); err != nil {
		return err
	}

	outBuf := []byte{}
	videoBuffer := data.Bytes()
	for offset := 0; offset < len(videoBuffer); {
		bufferLength := int(binary.BigEndian.Uint32(videoBuffer[offset : offset+headerLengthField]))
		if offset+bufferLength >= len(videoBuffer) {
			break
		}

		offset += headerLengthField
		outBuf = append(outBuf, []byte{0x00, 0x00, 0x00, 0x01}...)
		outBuf = append(outBuf, videoBuffer[offset:offset+bufferLength]...)

		offset += int(bufferLength)
	}

	h.videoRTPChan <- outBuf
	return nil
}

func (h *RTMPHandler) OnClose() {
	h.log.Infow("closing ingress RTMP session")

	if h.onClose != nil {
		h.onClose(h.resourceId)
	}
}
