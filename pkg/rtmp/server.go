package rtmp

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"path"
	"sync"

	"github.com/frostbyte73/core"
	"github.com/livekit/protocol/logger"
	protoutils "github.com/livekit/protocol/utils"
	log "github.com/sirupsen/logrus"
	"github.com/yutopp/go-flv"
	flvtag "github.com/yutopp/go-flv/tag"
	"github.com/yutopp/go-rtmp"
	rtmpmsg "github.com/yutopp/go-rtmp/message"

	"github.com/Harshitk-cp/rtmp_server/pkg/config"
	"github.com/Harshitk-cp/rtmp_server/pkg/errors"
	"github.com/Harshitk-cp/rtmp_server/pkg/params"
	"github.com/Harshitk-cp/rtmp_server/pkg/utils"
)

type RTMPServer struct {
	server   *rtmp.Server
	handlers sync.Map
}

func NewRTMPServer() *RTMPServer {
	return &RTMPServer{}
}

func (s *RTMPServer) Start(conf *config.Config, onPublish func(streamKey, resourceId string) (*params.Params, error)) error {
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
			if conf.Logging.JSON {
				l.SetFormatter(&log.JSONFormatter{})
			}
			lf := l.WithFields(conf.GetLoggerFields())

			log.WithFields(log.Fields{
				"remoteAddr": conn.RemoteAddr(),
			}).Info("Client connected")

			h := NewRTMPHandler()
			h.OnPublishCallback(func(streamKey, resourceId string) (*params.Params, error) {
				var params *params.Params
				var err error
				if onPublish != nil {
					params, err = onPublish(streamKey, resourceId)
					if err != nil {
						return nil, err
					}
				}

				s.handlers.Store(resourceId, h)

				go func() {
					err := h.RelayStream(h.RelayUrl, streamKey)
					if err != nil {
						log.Printf("Failed to relay stream: %v", err)
					}
				}()

				return params, nil
			})
			h.OnCloseCallback(func(resourceId string) {
				s.handlers.Delete(resourceId)
			})

			return conn, &rtmp.ConnConfig{
				Handler: h,

				ControlState: rtmp.StreamControlStateConfig{
					DefaultBandwidthWindowSize: 6 * 1024 * 1024 / 8,
				},

				Logger: lf,
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

func (s *RTMPServer) AssociateRelay(resourceId string, w io.WriteCloser) error {
	h, ok := s.handlers.Load(resourceId)
	if ok && h != nil {

		err := h.(*RTMPHandler).SetWriter(w)
		if err != nil {
			return err
		}
	} else {
		return errors.ErrIngressNotFound
	}

	return nil
}

func (s *RTMPServer) DissociateRelay(resourceId string) error {
	h, ok := s.handlers.Load(resourceId)
	if ok && h != nil {
		err := h.(*RTMPHandler).SetWriter(nil)
		if err != nil {
			return err
		}
	} else {
		return errors.ErrIngressNotFound
	}

	return nil
}

func (s *RTMPServer) CloseHandler(resourceId string) {
	h, ok := s.handlers.Load(resourceId)
	if ok && h != nil {
		h.(*RTMPHandler).Close()
	}
}

func (s *RTMPServer) Stop() error {
	return s.server.Close()
}

type RTMPHandler struct {
	rtmp.DefaultHandler

	flvEnc        *flv.Encoder
	params        *params.Params
	resourceId    string
	videoInit     *flvtag.VideoData
	audioInit     *flvtag.AudioData
	keyFrameFound bool
	mediaBuffer   *utils.PrerollBuffer
	RelayUrl      string

	log    logger.Logger
	closed core.Fuse

	onPublish func(streamKey, resourceId string) (*params.Params, error)
	onClose   func(resourceId string)
}

func NewRTMPHandler() *RTMPHandler {
	h := &RTMPHandler{
		log: logger.GetLogger(),
	}

	h.mediaBuffer = utils.NewPrerollBuffer(func() error {
		h.log.Infow("preroll buffer reset event")
		h.flvEnc = nil

		return nil
	})

	return h
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

const (
	chunkSize = 128
)

func (h *RTMPHandler) RelayStream(rtmpURL, streamKey string) error {
	conn, err := rtmp.Dial("rtmp", "localhost:1935", &rtmp.ConnConfig{
		Logger: log.StandardLogger(),
	})
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := conn.Connect(nil); err != nil {
		return err
	}

	stream, err := conn.CreateStream(nil, chunkSize)
	if err != nil {
		log.Fatalf("Failed to create stream: Err=%+v", err)
	}
	defer stream.Close()

	var timestamp uint32
	for {
		var flvTag flvtag.FlvTag

		var msg rtmpmsg.Message
		switch flvTag.TagType {
		case flvtag.TagTypeAudio:
			msg = &rtmpmsg.AudioMessage{
				Payload: flvTag.Data.(*flvtag.AudioData).Data,
			}
		case flvtag.TagTypeVideo:
			msg = &rtmpmsg.VideoMessage{
				Payload: flvTag.Data.(*flvtag.VideoData).Data,
			}
		default:
			continue
		}

		if err := stream.Write(0, timestamp, msg); err != nil {
			return err
		}

		timestamp += flvTag.Timestamp
	}
}

func (h *RTMPHandler) OnSetDataFrame(timestamp uint32, data *rtmpmsg.NetStreamSetDataFrame) error {
	if h.closed.IsBroken() {
		return io.EOF
	}

	if h.flvEnc == nil {
		err := h.initFlvEncoder()
		if err != nil {
			return err
		}
	}
	r := bytes.NewReader(data.Payload)

	var script flvtag.ScriptData
	if err := flvtag.DecodeScriptData(r, &script); err != nil {
		h.log.Errorw("failed to decode script data", err)
		return nil
	}

	if err := h.flvEnc.Encode(&flvtag.FlvTag{
		TagType:   flvtag.TagTypeScriptData,
		Timestamp: timestamp,
		Data:      &script,
	}); err != nil {
		h.log.Warnw("failed to forward script data", err)
		return err
	}

	return nil
}

func (h *RTMPHandler) OnAudio(timestamp uint32, payload io.Reader) error {
	var audio flvtag.AudioData
	if err := flvtag.DecodeAudioData(payload, &audio); err != nil {
		return err
	}

	audioBuffer := new(bytes.Buffer)
	if _, err := io.Copy(audioBuffer, audio.Data); err != nil {
		return err
	}
	audio.Data = audioBuffer

	if err := h.flvEnc.Encode(&flvtag.FlvTag{
		TagType:   flvtag.TagTypeAudio,
		Timestamp: timestamp,
		Data:      &audio,
	}); err != nil {
		log.Printf("Failed to write audio: %+v", err)
	}
	return nil
}

func (h *RTMPHandler) OnVideo(timestamp uint32, payload io.Reader) error {
	var video flvtag.VideoData
	if err := flvtag.DecodeVideoData(payload, &video); err != nil {
		return err
	}

	flvBody := new(bytes.Buffer)
	if _, err := io.Copy(flvBody, video.Data); err != nil {
		return err
	}
	video.Data = flvBody

	if err := h.flvEnc.Encode(&flvtag.FlvTag{
		TagType:   flvtag.TagTypeVideo,
		Timestamp: timestamp,
		Data:      &video,
	}); err != nil {
		log.Printf("Failed to write video: Err = %+v", err)
	}

	return nil
}

func (h *RTMPHandler) OnClose() {
	h.log.Infow("closing ingress RTMP session")

	h.mediaBuffer.Close()

	if h.onClose != nil {
		h.onClose(h.resourceId)
	}
}

func (h *RTMPHandler) SetWriter(w io.WriteCloser) error {
	return h.mediaBuffer.SetWriter(w)
}

func (h *RTMPHandler) Close() {
	h.closed.Break()
}

func (h *RTMPHandler) initFlvEncoder() error {
	h.keyFrameFound = false
	enc, err := flv.NewEncoder(h.mediaBuffer, flv.FlagsAudio|flv.FlagsVideo)
	if err != nil {
		return err
	}
	h.flvEnc = enc

	if h.videoInit != nil {
		if err := h.flvEnc.Encode(&flvtag.FlvTag{
			TagType:   flvtag.TagTypeVideo,
			Timestamp: 0,
			Data:      copyVideoTag(h.videoInit),
		}); err != nil {
			return err
		}
	}
	if h.audioInit != nil {
		if err := h.flvEnc.Encode(&flvtag.FlvTag{
			TagType:   flvtag.TagTypeAudio,
			Timestamp: 0,
			Data:      copyAudioTag(h.audioInit),
		}); err != nil {
			return err
		}
	}

	return nil
}

func copyVideoTag(in *flvtag.VideoData) *flvtag.VideoData {
	ret := *in
	ret.Data = bytes.NewBuffer(in.Data.(*bytes.Buffer).Bytes())

	return &ret
}

func copyAudioTag(in *flvtag.AudioData) *flvtag.AudioData {
	ret := *in
	ret.Data = bytes.NewBuffer(in.Data.(*bytes.Buffer).Bytes())

	return &ret
}
