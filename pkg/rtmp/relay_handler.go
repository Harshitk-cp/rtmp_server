package rtmp

import (
	"github.com/livekit/protocol/logger"
	"github.com/yutopp/go-rtmp"
	rtmpmsg "github.com/yutopp/go-rtmp/message"
)

type RTMPRelayHandler struct {
	rtmpServer *RTMPServer
}

func NewRTMPRelayHandler(rtmpServer *RTMPServer) *RTMPRelayHandler {
	return &RTMPRelayHandler{
		rtmpServer: rtmpServer,
	}
}

func (h *RTMPRelayHandler) OnPublish(conn *rtmp.Conn, cmd *rtmpmsg.NetStreamPublish) error {
	resourceId := cmd.PublishingName
	log := logger.Logger(logger.GetLogger().WithValues("resourceID", resourceId))
	log.Infow("relaying ingress")

	// err := h.rtmpServer.AssociateRelay(resourceId, conn)
	// if err != nil {
	// 	return err
	// }

	defer func() {
		h.rtmpServer.DissociateRelay(resourceId)
	}()

	return nil
}

func (h *RTMPRelayHandler) OnPlay(conn *rtmp.Conn, cmd *rtmpmsg.NetStreamPlay) error {
	return nil
}
