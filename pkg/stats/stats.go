package stats

type LocalMediaStatsGatherer struct {
	audioStats *MediaTrackStatGatherer
	videoStats *MediaTrackStatGatherer
}

const (
	InputAudio = "input-audio"
	InputVideo = "input-video"
)

func NewLocalMediaStatsGatherer() *LocalMediaStatsGatherer {
	return &LocalMediaStatsGatherer{
		audioStats: &MediaTrackStatGatherer{},
		videoStats: &MediaTrackStatGatherer{},
	}
}

func (s *LocalMediaStatsGatherer) RegisterTrackStats(kind string) *MediaTrackStatGatherer {
	gatherer := &MediaTrackStatGatherer{}
	switch kind {
	case InputAudio:
		s.audioStats = gatherer
	case InputVideo:
		s.videoStats = gatherer
	}
	return gatherer
}

func (s *LocalMediaStatsGatherer) GetAudioStats() *MediaTrackStatGatherer {
	return s.audioStats
}

func (s *LocalMediaStatsGatherer) GetVideoStats() *MediaTrackStatGatherer {
	return s.videoStats
}

type MediaTrackStatGatherer struct {
	totalBytes int64
}

func (t *MediaTrackStatGatherer) MediaReceived(bytes int64) {
	t.totalBytes += bytes
}

func (t *MediaTrackStatGatherer) TotalBytes() int64 {
	return t.totalBytes
}
