package types

import (
	"context"

	"github.com/Harshitk-cp/rtmp_server/pkg/ipc"
)

type MediaStatsUpdater interface {
	UpdateMediaStats(ctx context.Context, stats *ipc.MediaStats) error
}

type MediaStatGatherer interface {
	GatherStats(ctx context.Context) (*ipc.MediaStats, error)
}
