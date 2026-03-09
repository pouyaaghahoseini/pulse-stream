package analytics

import "testing"

func TestNewProcessor_StartsEmpty(t *testing.T) {
	processor := NewProcessor()

	if processor.TotalEvents != 0 {
		t.Errorf("expected TotalEvents to be 0, got %d", processor.TotalEvents)
	}

	if len(processor.PlatformStats) != 0 {
		t.Errorf("expected PlatformStats to be empty, got %d entries", len(processor.PlatformStats))
	}
}

func TestProcessEvent_FirstEventForPlatform(t *testing.T) {
	processor := NewProcessor()

	stats := processor.ProcessEvent("twitter", 20)

	if processor.TotalEvents != 1 {
		t.Errorf("expected TotalEvents to be 1, got %d", processor.TotalEvents)
	}

	if stats.Platform != "twitter" {
		t.Errorf("expected platform to be twitter, got %s", stats.Platform)
	}

	if stats.TotalPosts != 1 {
		t.Errorf("expected TotalPosts to be 1, got %d", stats.TotalPosts)
	}

	if stats.TotalEngagement != 20 {
		t.Errorf("expected TotalEngagement to be 20, got %d", stats.TotalEngagement)
	}

	if stats.AverageScore != 20.0 {
		t.Errorf("expected AverageScore to be 20.0, got %f", stats.AverageScore)
	}
}

func TestProcessEvent_MultipleEventsSamePlatform(t *testing.T) {
	processor := NewProcessor()

	processor.ProcessEvent("twitter", 20)
	stats := processor.ProcessEvent("twitter", 40)

	if processor.TotalEvents != 2 {
		t.Errorf("expected TotalEvents to be 2, got %d", processor.TotalEvents)
	}

	if stats.TotalPosts != 2 {
		t.Errorf("expected TotalPosts to be 2, got %d", stats.TotalPosts)
	}

	if stats.TotalEngagement != 60 {
		t.Errorf("expected TotalEngagement to be 60, got %d", stats.TotalEngagement)
	}

	if stats.AverageScore != 30.0 {
		t.Errorf("expected AverageScore to be 30.0, got %f", stats.AverageScore)
	}
}

func TestProcessEvent_MultiplePlatforms(t *testing.T) {
	processor := NewProcessor()

	processor.ProcessEvent("twitter", 20)
	processor.ProcessEvent("instagram", 90)

	twitterStats := processor.PlatformStats["twitter"]
	instagramStats := processor.PlatformStats["instagram"]

	if processor.TotalEvents != 2 {
		t.Errorf("expected TotalEvents to be 2, got %d", processor.TotalEvents)
	}

	if twitterStats.TotalPosts != 1 {
		t.Errorf("expected twitter TotalPosts to be 1, got %d", twitterStats.TotalPosts)
	}

	if twitterStats.AverageScore != 20.0 {
		t.Errorf("expected twitter AverageScore to be 20.0, got %f", twitterStats.AverageScore)
	}

	if instagramStats.TotalPosts != 1 {
		t.Errorf("expected instagram TotalPosts to be 1, got %d", instagramStats.TotalPosts)
	}

	if instagramStats.AverageScore != 90.0 {
		t.Errorf("expected instagram AverageScore to be 90.0, got %f", instagramStats.AverageScore)
	}
}