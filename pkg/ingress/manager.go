package ingress

import (
	"fmt"
	"sync"

	"github.com/Harshitk-cp/rtmp_server/pkg/config"
	"github.com/Harshitk-cp/rtmp_server/pkg/room"
	"github.com/Harshitk-cp/rtmp_server/pkg/rtmp"
)

type IngressManager struct {
	config     *config.Config
	sfuServer  *rtmp.SFUServer
	ingressMap sync.Map
}

func NewIngressManager(config *config.Config, sfuServer *rtmp.SFUServer) *IngressManager {
	return &IngressManager{
		config:    config,
		sfuServer: sfuServer,
	}
}

func (m *IngressManager) CreateIngress(streamKey string, r *room.Room) (*Ingress, error) {
	if ingress, ok := m.ingressMap.Load(streamKey); ok {
		return ingress.(*Ingress), nil
	}

	ingress := NewIngress(m.config, m.sfuServer, streamKey, r)
	if ingress == nil {
		return nil, fmt.Errorf("failed to create ingress for stream key: %s", streamKey)
	}

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
