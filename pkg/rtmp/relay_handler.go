package rtmp

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/livekit/protocol/logger"
	"github.com/livekit/psrpc"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type RTMPRelayHandler struct {
	rtmpServer *RTMPServer
}

func NewRTMPRelayHandler(rtmpServer *RTMPServer) *RTMPRelayHandler {
	return &RTMPRelayHandler{
		rtmpServer: rtmpServer,
	}
}

func (h *RTMPRelayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Error("SPET 1")
	var err error
	defer func() {
		log.Error("SPET 2")
		var psrpcErr psrpc.Error
		switch {
		case errors.As(err, &psrpcErr):
			w.WriteHeader(psrpcErr.ToHttp())
		case err == nil:
		default:
			log.Errorf("Error serving HTTP request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	resourceId := strings.TrimLeft(r.URL.Path, "/rtmp/")
	log := logger.Logger(logger.GetLogger().WithValues("resourceID", resourceId))
	log.Infow("relaying ingress")

	pr, pw := io.Pipe()
	defer pw.Close()

	done := make(chan error)
	go func() {
		_, copyErr := io.Copy(w, pr)
		done <- copyErr
		close(done)
	}()

	err = h.rtmpServer.AssociateRelay(resourceId, pw)
	if err != nil {
		logrus.Errorf("Failed to associate relay: %v", err)
		return
	}
	defer h.rtmpServer.DissociateRelay(resourceId)

	notify := w.(http.CloseNotifier).CloseNotify()
	select {
	case <-notify:
		logrus.Infof("Client closed the connection")
		return
	case err = <-done:
		if err != nil {
			logrus.Errorf("Error copying data: %v", err)
		}
	}
}
