package sojourner

import "time"

type perfEventType string

const (
	PERF_START perfEventType = "start"
	PERF_END   perfEventType = "end"
)

type PerfEvent interface {
	GetTimestamp() time.Time
	GetName() string
	GetType() perfEventType
	GetSubtractedTime() time.Duration
	SubtractSelf(time.Duration)
}

func newPerfEvent(type_ perfEventType, name string) PerfEvent {
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
