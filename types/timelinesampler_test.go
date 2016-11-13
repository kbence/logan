package types

import (
	"testing"
	"time"
)

type timelineSamplerTestCase struct {
	date             time.Time
	increment        uint64
	expectedPosition int
}

var location *time.Location

func newTimelineSamplerTestCase(date time.Time, increment uint64, expectedPosition int) *timelineSamplerTestCase {
	return &timelineSamplerTestCase{date: date, increment: increment, expectedPosition: expectedPosition}
}

func executeTimelineSamplerTest(t *testing.T, testCase *timelineSamplerTestCase) {
	startTime := time.Date(2016, 11, 13, 6, 44, 48, 0, location)
	endTime := time.Date(2016, 11, 13, 7, 44, 48, 0, location)
	sampler := NewTimelineSampler(NewTimeInterval(startTime, endTime), 30)

	sampler.Inc(testCase.date, testCase.increment)

	for idx, sample := range sampler.Samples {
		if idx == testCase.expectedPosition {
			if sample != testCase.increment {
				t.Errorf("Sample value (%d) must be %d at position %d!", sample,
					testCase.increment, idx)
			}
		} else {
			if sample != 0 {
				t.Errorf("Sample value (%d) must be 0 at position %d!", sample, idx)
			}
		}
	}
}

func TestTimelineSamplerInc(t *testing.T) {
	location, _ = time.LoadLocation("Local")

	executeTimelineSamplerTest(t, newTimelineSamplerTestCase(
		time.Date(2016, 11, 13, 6, 44, 48, 0, location), 1, 0))
	executeTimelineSamplerTest(t, newTimelineSamplerTestCase(
		time.Date(2016, 11, 13, 7, 14, 47, 0, location), 2, 14))
	executeTimelineSamplerTest(t, newTimelineSamplerTestCase(
		time.Date(2016, 11, 13, 7, 14, 48, 0, location), 3, 15))
	executeTimelineSamplerTest(t, newTimelineSamplerTestCase(
		time.Date(2016, 11, 13, 7, 44, 47, 0, location), 4, 29))

	// It should drop sampler that are out of range
	executeTimelineSamplerTest(t, newTimelineSamplerTestCase(
		time.Date(2016, 11, 13, 7, 44, 48, 0, location), 0, 0))
	executeTimelineSamplerTest(t, newTimelineSamplerTestCase(
		time.Date(2016, 11, 13, 6, 44, 47, 0, location), 0, 0))
}
