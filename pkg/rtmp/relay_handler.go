package rtmp

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/livekit/protocol/logger"
	"github.com/livekit/psrpc"

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
	var err error
	defer func() {
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

	resourceId := strings.TrimPrefix(r.URL.Path, "/rtmp/")
	log := logger.Logger(logger.GetLogger().WithValues("resourceID", resourceId))

	pr, pw := io.Pipe()
	defer pw.Close()

	done := make(chan error)
	go func() {
		_, copyErr := io.Copy(w, pr)
		done <- copyErr
		close(done)
	}()

	// err = h.rtmpServer.AssociateRelay(resourceId, pw)
	// if err != nil {
	// 	log.Errorw("Failed to associate relay: %v", err)
	// 	return
	// }
	// defer h.rtmpServer.DissociateRelay(resourceId)

	ctx := r.Context()
	select {
	case <-ctx.Done():
		log.Infow("Client closed the connection")
		err = ctx.Err()
		return
	case err = <-done:
		if err != nil {
			log.Errorw("Error copying data: %v", err)
		}
	}
}
