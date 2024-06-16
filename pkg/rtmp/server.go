package rtmp

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"path"

	"github.com/livekit/protocol/logger"
	protoutils "github.com/livekit/protocol/utils"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	flvtag "github.com/yutopp/go-flv/tag"
	"github.com/yutopp/go-rtmp"
	rtmpmsg "github.com/yutopp/go-rtmp/message"

	"github.com/Harshitk-cp/rtmp_server/pkg/config"
	"github.com/Harshitk-cp/rtmp_server/pkg/errors"
	"github.com/Harshitk-cp/rtmp_server/pkg/params"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
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
			l := logrus.StandardLogger()

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

	}
}

func (s *RTMPServer) Stop() error {
	return s.server.Close()
}

type RTMPHandler struct {
	rtmp.DefaultHandler

	params     *params.Params
	resourceId string

	RelayUrl string

	log logger.Logger

	onPublish func(streamKey, resourceId string) (*params.Params, error)
	onClose   func(resourceId string)

	pc           *webrtc.PeerConnection
	videoTrack   *webrtc.TrackLocalStaticRTP
	audioTrack   *webrtc.TrackLocalStaticRTP
	videoRTPChan chan *rtp.Packet
	audioRTPChan chan *rtp.Packet
}

func NewRTMPHandler() *RTMPHandler {
	return &RTMPHandler{
		videoRTPChan: make(chan *rtp.Packet),
		audioRTPChan: make(chan *rtp.Packet),
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

	err := h.initWebRTC()
	if err != nil {
		return err
	}

	return nil
}

func (h *RTMPHandler) initWebRTC() error {
	var err error
	h.pc, err = webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		return err
	}

	h.videoTrack, err = webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264}, "video", "pion")
	if err != nil {
		return err
	}

	h.audioTrack, err = webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus}, "audio", "pion")
	if err != nil {
		return err
	}

	_, err = h.pc.AddTrack(h.videoTrack)
	if err != nil {
		return err
	}

	_, err = h.pc.AddTrack(h.audioTrack)
	if err != nil {
		return err
	}

	return nil
}

func (h *RTMPHandler) ReceiveRTMPStream() (*webrtc.TrackLocalStaticRTP, *webrtc.TrackLocalStaticRTP, error) {
	videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264}, "video", "pion")
	if err != nil {
		return nil, nil, err
	}

	audioTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus}, "audio", "pion")
	if err != nil {
		return nil, nil, err
	}

	h.videoTrack = videoTrack
	h.audioTrack = audioTrack

	return videoTrack, audioTrack, nil
}

func (h *RTMPHandler) OnAudio(timestamp uint32, payload io.Reader) error {
	// var audio flvtag.AudioData
	// if err := flvtag.DecodeAudioData(payload, &audio); err != nil {
	// 	return err
	// }

	// audioBuffer := new(bytes.Buffer)
	// if _, err := io.Copy(audioBuffer, audio.Data); err != nil {
	// 	return err
	// }
	// audio.Data = audioBuffer

	// rtpPacket := &rtp.Packet{}
	// err := rtpPacket.Unmarshal(audioBuffer.Bytes())
	// if err != nil {
	// 	logrus.Errorf("Failed to unmarshal audio data into RTP packet: %v", err)
	// 	return err
	// }

	// h.audioRTPChan <- rtpPacket
	return nil
}

func (h *RTMPHandler) OnVideo(timestamp uint32, payload io.Reader) error {
	var video flvtag.VideoData
	if err := flvtag.DecodeVideoData(payload, &video); err != nil {
		return err
	}

	videoBuffer := new(bytes.Buffer)
	if _, err := io.Copy(videoBuffer, video.Data); err != nil {
		return err
	}
	video.Data = videoBuffer

	rtpPacket := &rtp.Packet{}
	err := rtpPacket.Unmarshal(videoBuffer.Bytes())
	if err != nil {
		logrus.Errorf("Failed to unmarshal video data into RTP packet: %v", err)
		return err
	}

	h.videoRTPChan <- rtpPacket
	return nil
}

func (h *RTMPHandler) sendAudioToWebRTC(data []byte) {

	rtpPacket := &rtp.Packet{}

	err := rtpPacket.Unmarshal(data)
	if err != nil {
		log.Error("Failed to unmarshal data into RTP packet: ", err)
		return
	}

	err = h.audioTrack.WriteRTP(rtpPacket)
	if err != nil {
		log.Error("Failed to write RTP packet to audio track: ", err)
	}
}

func (h *RTMPHandler) sendVideoToWebRTC(data []byte) {

	rtpPacket := &rtp.Packet{}

	err := rtpPacket.Unmarshal(data)
	if err != nil {
		log.Error("Failed to unmarshal data into RTP packet: ", err)
		return
	}

	err = h.videoTrack.WriteRTP(rtpPacket)
	if err != nil {
		log.Error("Failed to write RTP packet to video track: ", err)
	}
}

func (h *RTMPHandler) OnClose() {
	h.log.Infow("closing ingress RTMP session")

	if h.pc != nil {
		h.pc.Close()
	}

	if h.onClose != nil {
		h.onClose(h.resourceId)
	}
}
