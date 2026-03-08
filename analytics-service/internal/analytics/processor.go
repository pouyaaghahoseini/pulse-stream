package analytics

type PlatformStats struct {
	Platform        string
	TotalPosts      int
	TotalEngagement int
	AverageScore    float64
}

type Processor struct {
	TotalEvents   int
	PlatformStats map[string]*PlatformStats
}

func NewProcessor() *Processor {
	return &Processor{
		TotalEvents:   0,
		PlatformStats: make(map[string]*PlatformStats),
	}
}

func (p *Processor) ProcessEvent(platform string, engagementScore int) *PlatformStats {
	p.TotalEvents++

	stats, exists := p.PlatformStats[platform]
	if !exists {
		stats = &PlatformStats{
			Platform: platform,
		}
		p.PlatformStats[platform] = stats
	}

	stats.TotalPosts++
	stats.TotalEngagement += engagementScore
	stats.AverageScore = float64(stats.TotalEngagement) / float64(stats.TotalPosts)

	return stats
}