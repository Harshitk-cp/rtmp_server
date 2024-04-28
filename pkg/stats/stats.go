package stats

type LocalMediaStatsGatherer struct {
	// Implement necessary methods for media stats gathering if needed
}

const (
	InputAudio = "input-audio"
	InputVideo = "input-video"
)

func (s *LocalMediaStatsGatherer) RegisterTrackStats(kind string) *MediaTrackStatGatherer {
	return &MediaTrackStatGatherer{}
}

type MediaTrackStatGatherer struct {
	// Implement necessary methods for track stats gathering if needed
}

func (t *MediaTrackStatGatherer) MediaReceived(bytes int64) {
	// Implement media received logic if needed
}
