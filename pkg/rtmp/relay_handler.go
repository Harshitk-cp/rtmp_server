package rtmp

import (
	"io"
	"net/http"
	"strings"

	"github.com/Harshitk-cp/rtmp_server/pkg/errors"
	"github.com/livekit/protocol/logger"
	"github.com/livekit/psrpc"
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
			// Nothing, we already responded
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	resourceId := strings.TrimLeft(r.URL.Path, "/rtmp/")
	token := r.URL.Query().Get("token")

	log := logger.Logger(logger.GetLogger().WithValues("resourceID", resourceId))
	log.Infow("relaying ingress")

	pr, pw := io.Pipe()
	done := make(chan error)

	go func() {
		_, err = io.Copy(w, pr)
		done <- err
		close(done)
	}()

	err = h.rtmpServer.AssociateRelay(resourceId, token, pw)
	if err != nil {
		return
	}
	defer func() {
		pw.Close()
		h.rtmpServer.DissociateRelay(resourceId)
	}()

	err = <-done
}