package ingress

import (
	"sync"

	"github.com/Harshitk-cp/rtmp_server/pkg/config"
	"github.com/Harshitk-cp/rtmp_server/pkg/room"
	"github.com/Harshitk-cp/rtmp_server/pkg/rtmp"
)

type IngressManager struct {
	config     *config.Config
	rtmpServer *rtmp.RTMPServer
	ingressMap sync.Map
}

func NewIngressManager(config *config.Config, rtmpServer *rtmp.RTMPServer) *IngressManager {
	return &IngressManager{
		config:     config,
		rtmpServer: rtmpServer,
	}
}

func (m *IngressManager) CreateIngress(streamKey string, r *room.Room, participant *room.Participant) (*Ingress, error) {

	if ingress, ok := m.ingressMap.Load(streamKey); ok {
		return ingress.(*Ingress), nil
	}

	ingress := NewIngress(m.config, m.rtmpServer, streamKey, r, participant)

	m.ingressMap.Store(streamKey, ingress)

	return ingress, nil
}

func (m *IngressManager) RemoveIngress(streamKey string) error {
	ingress, ok := m.ingressMap.Load(streamKey)
	if !ok {
		return nil
	}
	ingress.(*Ingress).Stop()

	m.ingressMap.Delete(streamKey)

	return nil
}
