package ingress

import (
	"fmt"
	"sync"

	"github.com/Harshitk-cp/rtmp_server/pkg/config"
	"github.com/Harshitk-cp/rtmp_server/pkg/room"
	"github.com/Harshitk-cp/rtmp_server/pkg/rtmp"
	protoutils "github.com/livekit/protocol/utils"
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

func (m *IngressManager) CreateIngress(r *room.Room) (*Ingress, error) {

	ingressId := protoutils.NewGuid(protoutils.IngressPrefix)
	streamKey := protoutils.NewGuid(protoutils.RTMPResourcePrefix)

	ingress := NewIngress(m.config, m.sfuServer, streamKey, r, ingressId)
	if ingress == nil {
		return nil, fmt.Errorf("failed to create ingress for stream key: %s", streamKey)
	}

	m.ingressMap.Store(ingressId, ingress)

	return ingress, nil
}

func (m *IngressManager) RemoveIngress(ingressId string) error {
	ingress, ok := m.ingressMap.Load(ingressId)
	if !ok {
		return nil
	}
	ingress.(*Ingress).Stop()

	m.ingressMap.Delete(ingressId)

	return nil
}
