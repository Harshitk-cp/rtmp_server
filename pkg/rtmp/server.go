package rtmp

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"sync"
	"time"

	"github.com/Harshitk-cp/rtmp_server/pkg/config"
	"github.com/Harshitk-cp/rtmp_server/pkg/errors"
	"github.com/Harshitk-cp/rtmp_server/pkg/params"
	"github.com/Harshitk-cp/rtmp_server/pkg/stats"
	"github.com/Harshitk-cp/rtmp_server/pkg/types"
	"github.com/Harshitk-cp/rtmp_server/pkg/utils"
	"github.com/frostbyte73/core"
	"github.com/livekit/protocol/logger"
	protoutils "github.com/livekit/protocol/utils"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/yutopp/go-flv"
	flvtag "github.com/yutopp/go-flv/tag"
	"github.com/yutopp/go-rtmp"
	rtmpmsg "github.com/yutopp/go-rtmp/message"
)

type RTMPServer struct {
	server   *rtmp.Server
	handlers sync.Map
}

func NewRTMPServer() *RTMPServer {
	return &RTMPServer{}
}

func (s *RTMPServer) Start(port int, conf *config.Config, onPublish func(streamKey, resourceId string) (*params.Params, *stats.LocalMediaStatsGatherer, error)) error {

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		logger.Errorw("failed to start TCP listener", err, "port", port)
		return err
	}

	h, _ := NewRTMPHandler()
	h.OnPublishCallback(func(streamKey, resourceId string) (*params.Params, *stats.LocalMediaStatsGatherer, error) {
		var params *params.Params
		var statsGatherer *stats.LocalMediaStatsGatherer
		var err error
		if onPublish != nil {
			params, statsGatherer, err = onPublish(streamKey, resourceId)
			if err != nil {
				return nil, nil, err
			}
		}

		s.handlers.Store(resourceId, h)

		return params, statsGatherer, nil
	})
	h.OnCloseCallback(func(resourceId string) {
		s.handlers.Delete(resourceId)
	})

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

	s.server = srv
	go func() {
		if err := srv.Serve(listener); err != nil {
			logger.Errorw("failed to start RTMP server", err)
		}
	}()

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
	trackStats    map[types.StreamKind]*stats.MediaTrackStatGatherer
	params        *params.Params
	resourceId    string
	videoInit     *flvtag.VideoData
	audioInit     *flvtag.AudioData
	keyFrameFound bool
	mediaBuffer   *utils.PrerollBuffer
	RelayUrl      string
	relayConn     *rtmp.ClientConn
	relayStream   *rtmp.Stream
	flvFile       *os.File

	log    logger.Logger
	closed core.Fuse

	onPublish func(streamKey, resourceId string) (*params.Params, *stats.LocalMediaStatsGatherer, error)
	onClose   func(resourceId string)
}

func NewRTMPHandler() (*RTMPHandler, *utils.PrerollBuffer) {
	h := &RTMPHandler{
		log:        logger.GetLogger(),
		trackStats: make(map[types.StreamKind]*stats.MediaTrackStatGatherer),
	}

	h.mediaBuffer = utils.NewPrerollBuffer(func() error {
		h.log.Infow("preroll buffer reset event")
		h.flvEnc = nil

		return nil
	})

	return h, h.mediaBuffer
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

func (h *RTMPHandler) SetWriter(w io.WriteCloser) error {
	return h.mediaBuffer.SetWriter(w)
}

func (h *RTMPHandler) OnPublishCallback(cb func(streamKey, resourceId string) (*params.Params, *stats.LocalMediaStatsGatherer, error)) {
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
		params, statsGatherer, err := h.onPublish(streamKey, h.resourceId)
		if err != nil {
			return err
		}
		h.params = params

		if statsGatherer == nil {
			statsGatherer = stats.NewLocalMediaStatsGatherer()
		}

		h.trackStats[types.Audio] = statsGatherer.RegisterTrackStats(stats.InputAudio)
		h.trackStats[types.Video] = statsGatherer.RegisterTrackStats(stats.InputVideo)
	}

	logrus.WithFields(logrus.Fields{
		"appName":        appName,
		"streamKey":      streamKey,
		"resourceID":     h.resourceId,
		"publishingName": cmd.PublishingName,
	}).Info("Received publish request")

	h.log.Infow("Received a new published stream")

	// p := filepath.Join(
	// 	os.TempDir(),
	// 	filepath.Clean(filepath.Join("../", fmt.Sprintf("%s.flv", cmd.PublishingName))),
	// )
	// f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY, 0666)
	// if err != nil {
	// 	return err
	// }
	// h.flvFile = f

	// enc, err := flv.NewEncoder(f, flv.FlagsAudio|flv.FlagsVideo)
	// if err != nil {
	// 	_ = f.Close()
	// 	return err
	// }
	// h.flvEnc = enc

	if h.relayStream == nil {

		go func() {
			err := h.RelayStream("key")
			logrus.Info("created new relay stream")

			if err != nil {
				logrus.Printf("Failed to relay stream: %v", err)
			}
		}()
	} else {
		logrus.Info("Relay stream already created")
	}

	return nil
}

const (
	chunkSize = 8192
)

func (s *RTMPHandler) RelayStream(streamKey string) error {

	dialer := &net.Dialer{
		Timeout:   180 * time.Second,
		KeepAlive: 180 * time.Second,
	}

	conn, err := rtmp.DialWithDialer(dialer, "rtmp", "127.0.0.1:9090", &rtmp.ConnConfig{
		Logger: logrus.StandardLogger(),
	})
	if err != nil {
		return err
	}
	defer conn.Close()

	connectCommand := &rtmpmsg.NetConnectionConnect{
		Command: rtmpmsg.NetConnectionConnectCommand{
			App:            "streamix",
			Type:           "nonprivate",
			FlashVer:       "FMLE/3.0",
			TCURL:          "rtmp://localhost",
			Fpad:           false,
			Capabilities:   239,
			AudioCodecs:    3575,
			VideoCodecs:    252,
			VideoFunction:  1,
			ObjectEncoding: 0,
		},
	}

	if err := conn.Connect(connectCommand); err != nil {
		return err
	}

	createStreamCommand := &rtmpmsg.NetConnectionCreateStream{}

	stream, err := conn.CreateStream(createStreamCommand, chunkSize)
	if err != nil {
		return err
	}
	defer stream.Close()

	if err := stream.Publish(&rtmpmsg.NetStreamPublish{
		PublishingName: "streamix",
		PublishingType: "live",
	}); err != nil {
		logrus.Errorf("Failed to publish stream: %v", err)
	}
	s.relayStream = stream

	return nil
}

func (h *RTMPHandler) OnSetDataFrame(timestamp uint32, data *rtmpmsg.NetStreamSetDataFrame) error {
	r := bytes.NewReader(data.Payload)

	if h.flvEnc == nil {
		if err := h.initFlvEncoder(); err != nil {
			return err
		}
	}

	var script flvtag.ScriptData
	if err := flvtag.DecodeScriptData(r, &script); err != nil {
		log.Printf("Failed to decode script data: Err = %+v", err)
		return nil
	}

	log.Printf("SetDataFrame: Script = %#v", script)

	if err := h.flvEnc.Encode(&flvtag.FlvTag{
		TagType:   flvtag.TagTypeScriptData,
		Timestamp: timestamp,
		Data:      &script,
	}); err != nil {
		log.Printf("Failed to write script data: Err = %+v", err)
	}

	return nil
}

func (h *RTMPHandler) OnAudio(timestamp uint32, payload io.Reader) error {
	var audio flvtag.AudioData
	if err := flvtag.DecodeAudioData(payload, &audio); err != nil {
		return err
	}

	flvBody := new(bytes.Buffer)
	if _, err := io.Copy(flvBody, audio.Data); err != nil {
		return err
	}
	audio.Data = flvBody

	// log.Printf("FLV Audio Data: Timestamp = %d, SoundFormat = %+v, SoundRate = %+v, SoundSize = %+v, SoundType = %+v, AACPacketType = %+v, Data length = %+v",
	// 	timestamp,
	// 	audio.SoundFormat,
	// 	audio.SoundRate,
	// 	audio.SoundSize,
	// 	audio.SoundType,
	// 	audio.AACPacketType,
	// 	len(flvBody.Bytes()),
	// )

	if err := h.flvEnc.Encode(&flvtag.FlvTag{
		TagType:   flvtag.TagTypeAudio,
		Timestamp: timestamp,
		Data:      &audio,
	}); err != nil {
		log.Printf("Failed to write audio: Err = %+v", err)
	}

	if h.relayStream != nil {
		// msg := &rtmpmsg.AudioMessage{
		// 	Payload: audio.Data,
		// }
		// err := h.writeToStream(4, timestamp, msg)
		// if err != nil {
		// 	logrus.Errorf("Failed to write audio to relay stream: %v", err)
		// }
	} else {
		logrus.Error("Relay stream is not initialized on audio")
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

	// log.Printf("FLV Video Data: Timestamp = %d, FrameType = %+v, CodecID = %+v, AVCPacketType = %+v, CT = %+v, Data length = %+v",
	// 	timestamp,
	// 	video.FrameType,
	// 	video.CodecID,
	// 	video.AVCPacketType,
	// 	video.CompositionTime,
	// 	len(flvBody.Bytes()),
	// )

	if err := h.flvEnc.Encode(&flvtag.FlvTag{
		TagType:   flvtag.TagTypeVideo,
		Timestamp: timestamp,
		Data:      &video,
	}); err != nil {
		log.Printf("Failed to write video: Err = %+v", err)
	}

	if h.relayStream != nil {
		msg := &rtmpmsg.VideoMessage{
			Payload: video.Data,
		}

		err := h.writeToStream(6, timestamp, msg)
		if err != nil {
			logrus.Errorf("Failed to write video to relay stream: %v", err)
		}
	} else {
		logrus.Error("Relay stream is not initialized on video")
	}

	return nil
}

func (h *RTMPHandler) writeToStream(chunkStreamID int, timestamp uint32, msg rtmpmsg.Message) error {

	stream := h.relayStream
	if stream == nil {
		log.Error("no relay stream found")
		return nil
	}
	logrus.Infof("Message payload length before encoding: %d", msg)

	chunkSize := 1024
	payload := new(bytes.Buffer)
	enc := rtmpmsg.NewEncoder(payload)
	if err := enc.Encode(msg); err != nil {
		return err
	}

	data := payload.Bytes()

	chunk := data
	if len(chunk) > chunkSize {
		chunk = chunk[:chunkSize]
	}

	logrus.Infof("setting chunk : ChunkStreamID = %d, Timestamp = %d, Length = %d", chunkStreamID, timestamp, len(chunk))
	if err := stream.WriteSetChunkSize(uint32(chunkSize)); err != nil {
		return err
	}

	logrus.Infof("Writing data to relay stream: ChunkStreamID = %d, Timestamp = %d, Length = %d", chunkStreamID, timestamp, len(chunk))

	if err := stream.Write(chunkStreamID, timestamp, msg); err != nil {
		logrus.Errorf("Failed to write data to relay stream: %v", err)
		return err
	}
	logrus.Infof("Writing data to relay stream: ChunkStreamID = %d, Timestamp = %d, Length = %d", chunkStreamID, timestamp, len(chunk))

	errCh := make(chan error, 1)
	go func() {
		errCh <- stream.Write(chunkStreamID, timestamp, nil)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			logrus.Errorf("Failed to write data to relay stream: %v", err)
			return err
		}
	case <-time.After(10 * time.Second):
		log.Error("timeout while writing to relay stream")
		return nil
	}

	logrus.Info("Data written to relay stream successfully")

	return nil
}

func (h *RTMPHandler) OnClose() {
	h.log.Infow("closing ingress RTMP session")

	h.mediaBuffer.Close()

	if h.relayStream != nil {
		h.relayStream.Close()
		h.relayStream = nil
	}
	if h.relayConn != nil {
		h.relayConn.Close()
		h.relayConn = nil
	}

	if h.onClose != nil {
		h.onClose(h.resourceId)
	}
}

func (h *RTMPHandler) Close() {
	h.closed.Break()
}

func (h *RTMPHandler) initFlvEncoder() error {
	h.keyFrameFound = false

	_, err := h.mediaBuffer.Write([]byte{0x46, 0x4C, 0x56})
	if err != nil {
		return err
	}

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
