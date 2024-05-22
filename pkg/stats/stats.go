package stats

type LocalMediaStatsGatherer struct {
}

const (
	InputAudio = "input-audio"
	InputVideo = "input-video"
)

func (s *LocalMediaStatsGatherer) RegisterTrackStats(kind string) *MediaTrackStatGatherer {
	return &MediaTrackStatGatherer{}
}

type MediaTrackStatGatherer struct {
}

func (t *MediaTrackStatGatherer) MediaReceived(bytes int64) {
}
