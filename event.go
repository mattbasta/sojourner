package sojourner

import "time"

type perfEventType string

const (
	PERF_START perfEventType = "start"
	PERF_END   perfEventType = "end"
)

type perfEvent interface {
	// Returns the timestamp that the event began taking place at.
	GetTimestamp() time.Time
	// Returns the name of the event.
	GetName() string
	// Returns the type. This should be either PERF_START or PERF_END.
	GetType() perfEventType
	// Returns the amount of time that sub-tasks consumed.
	GetSubtractedTime() time.Duration
	// Adds the passed duration to the amount of time subtracted from the
	// cumulative duration ("self" time).
	SubtractSelf(time.Duration)

	canComplete(perfEvent) bool
}

func newPerfEvent(type_ perfEventType, name string) perfEvent {
	now := time.Now()

	switch type_ {
	case PERF_START:
		return &(startPerfEvent{now, name, 0})
	case PERF_END:
		return &(endPerfEvent{now, name, 0})
	}
	return nil
}

type startPerfEvent struct {
	timestamp      time.Time
	name           string
	subtractedTime time.Duration
}

func (self startPerfEvent) GetTimestamp() time.Time          { return self.timestamp }
func (self startPerfEvent) GetName() string                  { return self.name }
func (self startPerfEvent) GetType() perfEventType           { return PERF_START }
func (self startPerfEvent) GetSubtractedTime() time.Duration { return self.subtractedTime }
func (self *startPerfEvent) SubtractSelf(dur time.Duration)  { self.subtractedTime += dur }
func (self startPerfEvent) canComplete(e perfEvent) bool     { return false }

type endPerfEvent struct {
	timestamp      time.Time
	name           string
	subtractedTime time.Duration
}

func (self endPerfEvent) GetTimestamp() time.Time          { return self.timestamp }
func (self endPerfEvent) GetName() string                  { return self.name }
func (self endPerfEvent) GetType() perfEventType           { return PERF_END }
func (self endPerfEvent) GetSubtractedTime() time.Duration { return self.subtractedTime }
func (self *endPerfEvent) SubtractSelf(dur time.Duration)  { self.subtractedTime += dur }
func (self endPerfEvent) canComplete(e perfEvent) bool     { return e.GetName() == self.name }
