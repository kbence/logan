package types

import "time"

type TimelineSampler struct {
	Interval TimeInterval
	Samples  []uint64
}

func NewTimelineSampler(interval *TimeInterval, samples int) *TimelineSampler {
	return &TimelineSampler{Interval: *interval, Samples: make([]uint64, samples)}
}

func (s *TimelineSampler) Inc(date time.Time, increment uint64) {
	startNano := s.Interval.StartTime.UnixNano()
	endNano := s.Interval.EndTime.UnixNano()
	dateNano := date.UnixNano()

	if dateNano < startNano || dateNano >= endNano {
		return
	}

	diff := endNano - startNano
	pos := dateNano - startNano

	s.Samples[(pos*int64(len(s.Samples)))/diff] += increment
}
