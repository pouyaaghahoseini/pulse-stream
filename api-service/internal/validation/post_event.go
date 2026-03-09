package validation

import "errors"

type PostEvent struct {
	PostID          string
	Platform        string
	ContentType     string
	EngagementScore int
}

func ValidatePostEvent(event PostEvent) error {
	if event.PostID == "" || event.Platform == "" || event.ContentType == "" {
		return errors.New("missing required fields")
	}

	if event.EngagementScore < 0 {
		return errors.New("engagement_score cannot be negative")
	}

	return nil
}