package utils

import "time"

type UpdateTimer struct {
	lastUpdateTime *time.Time
	lastUpdate     uint64
	received       uint64
	nextUpdate     uint64
}

func NewUpdateTimer() *UpdateTimer {
	return &UpdateTimer{}
}

func (t *UpdateTimer) IsUpdateNeeded() bool {
	t.received++

	if t.lastUpdateTime == nil || t.received >= t.nextUpdate {
		currentTime := time.Now()

		if t.lastUpdateTime == nil {
			t.nextUpdate = t.received + 10
		} else {
			next := 1000000000 * (t.received - t.lastUpdate) /
				uint64(currentTime.Sub(*t.lastUpdateTime))
			t.nextUpdate = t.received + next
		}

		t.lastUpdateTime = &currentTime
		t.lastUpdate = t.received

		return true
	}

	return false
}
