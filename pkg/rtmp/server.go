package rtmp

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"path"
	"sync"

	"github.com/Glimesh/go-fdkaac/fdkaac"
	"github.com/livekit/protocol/logger"
	protoutils "github.com/livekit/protocol/utils"
	log "github.com/sirupsen/logrus"
	flvtag "github.com/yutopp/go-flv/tag"
	"github.com/yutopp/go-rtmp"
	rtmpmsg "github.com/yutopp/go-rtmp/message"
	"gopkg.in/hraban/opus.v2"

	"github.com/Harshitk-cp/rtmp_server/pkg/config"
	"github.com/Harshitk-cp/rtmp_server/pkg/errors"
	"github.com/Harshitk-cp/rtmp_server/pkg/room"

	"github.com/Harshitk-cp/rtmp_server/pkg/params"
	"github.com/Harshitk-cp/rtmp_server/pkg/webhook"
)

type RTMPServer struct {
	server   *rtmp.Server
	sfu      *SFUServer
	handlers map[string]*RTMPHandler
	mu       sync.Mutex
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
			s.mu.Lock()
			s.handlers[h.resourceId] = h
			s.mu.Unlock()

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

func (h *RTMPHandler) OnConnect(timestamp uint32, cmd *rtmpmsg.NetConnectionConnect) error {
	log.Printf("OnConnect: %#v", cmd)
	h.audioClockRate = 48000
	return nil
}

func (s *RTMPServer) CloseHandler(resourceId string) {
	s.mu.Lock()
	defer s.mu.Unlock()
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

	audioDecoder   *fdkaac.AacDecoder
	audioEncoder   *opus.Encoder
	audioBuffer    []byte
	audioClockRate uint32

	videoBuffer    []byte
	pcm16Buffer    []int16
	opusDataBuffer []byte

	mu sync.Mutex

	webhookManager *webhook.WebhookManager
	roomManager    *room.RoomManager
}

func NewRTMPHandler(webhookManager *webhook.WebhookManager, roomManager *room.RoomManager) *RTMPHandler {
	return &RTMPHandler{
		videoRTPChan:   make(chan []byte, 100),
		audioRTPChan:   make(chan []byte, 100),
		pcm16Buffer:    make([]int16, 960*2), // Allocate once and reuse
		opusDataBuffer: make([]byte, 4000),   // Allocate once and reuse
		webhookManager: webhookManager,
		roomManager:    roomManager,
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

func (h *RTMPHandler) SetOpusCtl() {
	h.audioEncoder.SetMaxBandwidth(opus.Bandwidth(opus.Fullband))
	h.audioEncoder.SetComplexity(9)
	h.audioEncoder.SetBitrateToMax()
	// h.audioEncoder.SetBitrate(5000)
	h.audioEncoder.SetInBandFEC(true)
}

func (h *RTMPHandler) initAudio() error {
	encoder, err := opus.NewEncoder(48000, 2, opus.AppAudio)
	if err != nil {
		println(err.Error())
		return err
	}
	h.audioEncoder = encoder
	h.SetOpusCtl()
	h.audioDecoder = fdkaac.NewAacDecoder()

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
	if data.Len() <= 0 {
		log.Println("no audio data", timestamp, payload)
		return fmt.Errorf("no audio data")
	}
	datas := data.Bytes()

	if audio.AACPacketType == flvtag.AACPacketTypeSequenceHeader {
		log.Println("Created new codec ", hex.EncodeToString(datas))
		if err := h.initAudio(); err != nil {
			log.Println("error initializing Audio", err)
			return fmt.Errorf("can't initialize codec: %s", err.Error())
		}
		if err := h.audioDecoder.InitRaw(datas); err != nil {
			log.Println("error initializing stream", err)
			return fmt.Errorf("can't initialize codec with %s", hex.EncodeToString(datas))
		}
		return nil
	}

	pcm, err := h.audioDecoder.Decode(datas)
	if err != nil {
		log.Println("decode error: ", hex.EncodeToString(datas), err)
		return fmt.Errorf("decode error")
	}

	blockSize := 960
	h.mu.Lock()
	h.audioBuffer = append(h.audioBuffer, pcm...)
	for len(h.audioBuffer) >= blockSize*4 {
		for i := 0; i < blockSize*2; i++ {
			h.pcm16Buffer[i] = int16(binary.LittleEndian.Uint16(h.audioBuffer[i*2:]))
		}

		n, err := h.audioEncoder.Encode(h.pcm16Buffer, h.opusDataBuffer)
		if err != nil {
			h.mu.Unlock()
			return err
		}
		h.audioRTPChan <- h.opusDataBuffer[:n]
		h.audioBuffer = h.audioBuffer[blockSize*4:]
	}

	h.mu.Unlock()
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

	videoBuffer := data.Bytes()
	outBuf := h.videoBuffer[:0]
	for offset := 0; offset < len(videoBuffer); {
		if offset+headerLengthField >= len(videoBuffer) {
			break
		}
		bufferLength := int(binary.BigEndian.Uint32(videoBuffer[offset : offset+headerLengthField]))
		offset += headerLengthField
		if offset+bufferLength > len(videoBuffer) {
			break
		}
		outBuf = append(outBuf, []byte{0x00, 0x00, 0x00, 0x01}...)
		outBuf = append(outBuf, videoBuffer[offset:offset+bufferLength]...)
		offset += bufferLength
	}
	h.videoRTPChan <- outBuf
	return nil
}

func (h *RTMPHandler) OnClose() {
	h.log.Infow("closing ingress RTMP session")

	// h.webhookManager.SendWebhook("streamKey", "ingress_ended", "IN_sAw3KXXai3Ho")

	if h.onClose != nil {
		h.onClose(h.resourceId)
	}
}
