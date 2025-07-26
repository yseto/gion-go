package revisit

import (
	"math"
	"time"
)

type Manager struct {
	now    time.Time
	factor float64

	minRevisitDuration, maxRevisitDuration time.Duration
}

func New(now time.Time, factor float64, minRevisitDuration, maxRevisitDuration time.Duration) *Manager {
	return &Manager{
		now:                now,
		factor:             factor,
		minRevisitDuration: minRevisitDuration,
		maxRevisitDuration: maxRevisitDuration,
	}
}

func (i *Manager) revisitDuration(latestEntryAt time.Time) time.Duration {
	f := i.now.Sub(latestEntryAt) / time.Duration(i.minRevisitDuration)

	duration := time.Duration(math.Sqrt(float64(f)*i.factor)) * i.minRevisitDuration
	return max(min(duration, i.maxRevisitDuration), i.minRevisitDuration)
}

func (i *Manager) RevisitAt(latestEntryAt time.Time) time.Time {
	return i.now.Add(i.revisitDuration(latestEntryAt))
}

func (i *Manager) RevisitAtMaxRevisitDuration() time.Time {
	return i.now.Add(i.maxRevisitDuration)
}
