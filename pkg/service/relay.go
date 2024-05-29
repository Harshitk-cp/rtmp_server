package service

import (
	"fmt"
	"net/http"

	"github.com/Harshitk-cp/rtmp_server/pkg/config"
	"github.com/Harshitk-cp/rtmp_server/pkg/rtmp"
	"github.com/livekit/protocol/logger"
)

type Relay struct {
	server     *http.Server
	rtmpServer *rtmp.RTMPServer
}

func NewRelay(rtmpServer *rtmp.RTMPServer) *Relay {
	return &Relay{
		rtmpServer: rtmpServer,
	}
}

func (r *Relay) Start(conf *config.Config) error {
	port := conf.HTTPRelayPort

	mux := http.NewServeMux()

	if r.rtmpServer != nil {
		h := rtmp.NewRTMPRelayHandler(r.rtmpServer)
		mux.Handle("/rtmp/", h)
	}

	r.server = &http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf("localhost:%d", port),
	}

	go func() {
		err := r.server.ListenAndServe()
		logger.Debugw("Relay stopped", "error", err)
	}()

	return nil
}

func (r *Relay) Stop() error {
	return r.server.Close()
}
