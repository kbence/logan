package utils

import (
	"testing"
	"time"

	"github.com/kbence/logan/types"
)

func testBrokenComponents(t *testing.T, spec string, expected []string) {
	parts := breakTimeIntervalSpec(spec)

	if len(parts) != len(expected) {
		t.Errorf("len(parts) (%d) != len(expected) (%d)", len(parts), len(expected))
	}

	for i, val := range parts {
		if val != expected[i] {
			t.Errorf("parts[%d] (%s) != expected[%d] (%s)", i, val, i, expected[i])
		}
	}
}

func TestBreakTimeIntervalSpecSeparatesProperly(t *testing.T) {
	testBrokenComponents(t, "12:30", []string{"12:30", "", ""})
	testBrokenComponents(t, "-1h", []string{"", "1h", ""})
	testBrokenComponents(t, "-1h+5m", []string{"", "1h", "5m"})
	testBrokenComponents(t, "12:30+5m", []string{"12:30", "", "5m"})
	testBrokenComponents(t, "12:30-30m+5m", []string{"12:30", "30m", "5m"})
}

func testTimeReferenceMatch(t *testing.T, ref string, expected time.Time) {
	timeRef := parseTimeReference(ref, time.Date(2016, 11, 3, 13, 5, 12, 0, expected.Location()))

	if timeRef != expected {
		t.Errorf("Returned timeRef (%s) != expected (%s)", timeRef.String(), expected.String())
	}
}

func TestParseTimeReferenceIfOnlyHourIsSpecified(t *testing.T) {
	location, _ := time.LoadLocation("Local")

	testTimeReferenceMatch(t, "12:30", time.Date(2016, 11, 3, 12, 30, 0, 0, location))
	testTimeReferenceMatch(t, "12:30:54", time.Date(2016, 11, 3, 12, 30, 54, 0, location))
	testTimeReferenceMatch(t, "13:30:54", time.Date(2016, 11, 2, 13, 30, 54, 0, location))
}

func testTimeIntervalMatch(t *testing.T, spec string, expected *types.TimeInterval) {
	interval := ParseTimeInterval(spec, time.Date(2016, 11, 3, 13, 5, 12, 1234,
		expected.StartTime.Location()))

	if interval.StartTime != expected.StartTime {
		t.Errorf("interval.StartTime (%s) != expected.StartTime (%s)",
			interval.StartTime, expected.StartTime)
	}

	if interval.EndTime != expected.EndTime {
		t.Errorf("interval.EndTime (%s) != expected.EndTime (%s)",
			interval.EndTime, expected.EndTime)
	}
}

func TestParseTimeInterval(t *testing.T) {
	location, _ := time.LoadLocation("Local")

	testTimeIntervalMatch(t, "-1h",
		types.NewTimeInterval(
			time.Date(2016, 11, 3, 12, 5, 12, 1234, location),
			time.Date(2016, 11, 3, 13, 5, 12, 1234, location)))

	testTimeIntervalMatch(t, "-1h+5m8s",
		types.NewTimeInterval(
			time.Date(2016, 11, 3, 12, 5, 12, 1234, location),
			time.Date(2016, 11, 3, 12, 10, 20, 1234, location)))

	testTimeIntervalMatch(t, "11:30+5m",
		types.NewTimeInterval(
			time.Date(2016, 11, 3, 11, 30, 0, 0, location),
			time.Date(2016, 11, 3, 11, 35, 0, 0, location)))

	testTimeIntervalMatch(t, "15:30+5m",
		types.NewTimeInterval(
			time.Date(2016, 11, 2, 15, 30, 0, 0, location),
			time.Date(2016, 11, 2, 15, 35, 0, 0, location)))

	testTimeIntervalMatch(t, "12:30-48h+15m",
		types.NewTimeInterval(
			time.Date(2016, 11, 1, 12, 30, 0, 0, location),
			time.Date(2016, 11, 1, 12, 45, 0, 0, location)))
}
