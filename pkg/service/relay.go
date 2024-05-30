package service

import (
	"github.com/Harshitk-cp/rtmp_server/pkg/rtmp"
	"github.com/go-chi/chi/v5"
)

type Relay struct {
	rtmpServer *rtmp.RTMPServer
	router     *chi.Mux
}

func NewRelay(rtmpServer *rtmp.RTMPServer, router *chi.Mux) *Relay {
	return &Relay{
		rtmpServer: rtmpServer,
		router:     router,
	}
}

func (r *Relay) Start() error {
	r.router.Handle("/rtmp/*", rtmp.NewRTMPRelayHandler(r.rtmpServer))
	return nil
}
