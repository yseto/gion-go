package revisit

import (
	"testing"
	"time"
)

func TestRevisit(t *testing.T) {
	ts := New(time.Date(2025, 1, 10, 0, 30, 0, 0, time.UTC), 30, 30*time.Minute, 72*time.Hour)

	actual := ts.revisitDuration(time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC))
	if actual != 150*time.Minute {
		t.Error("invalid")
	}
}

func TestRevisitLoop(t *testing.T) {
	nowTs := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)
	revisit := nowTs

	maxRevisitDuration := 7 * 24 * time.Hour

	latestEntryAt := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)

	var duration time.Duration
	for duration < maxRevisitDuration {
		nowTs = revisit

		ts := New(nowTs, 30, 30*time.Minute, maxRevisitDuration)

		revisit = ts.RevisitAt(latestEntryAt)
		duration = ts.revisitDuration(latestEntryAt)
		t.Log(revisit, duration)
	}
}
