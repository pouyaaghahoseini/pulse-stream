package validation

import "testing"

func TestValidatePostEvent_ValidEvent(t *testing.T) {
	event := PostEvent{
		PostID:          "p_001",
		Platform:        "twitter",
		ContentType:     "text",
		EngagementScore: 42,
	}

	err := ValidatePostEvent(event)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidatePostEvent_MissingRequiredFields(t *testing.T) {
	event := PostEvent{
		PostID:          "",
		Platform:        "twitter",
		ContentType:     "text",
		EngagementScore: 42,
	}

	err := ValidatePostEvent(event)

	if err == nil {
		t.Errorf("expected error for missing required fields, got nil")
	}
}

func TestValidatePostEvent_NegativeEngagementScore(t *testing.T) {
	event := PostEvent{
		PostID:          "p_001",
		Platform:        "twitter",
		ContentType:     "text",
		EngagementScore: -1,
	}

	err := ValidatePostEvent(event)

	if err == nil {
		t.Errorf("expected error for negative engagement score, got nil")
	}
}