package sojourner

import (
	"sync"
	"time"
)

// The main performance monitoring object. Each monitor represents a single
// set of serial (though not necessarily synchronous) instrumented tasks.
type Monitor struct {
	created     time.Time
	inboundData chan perfEvent
	perfStack   []perfEvent

	aggregatedDataCumulative map[string]time.Duration
	aggregatedDataSelf       map[string]time.Duration

	aggregateLock sync.RWMutex
}

// Returns a new monitor object.
//
// You should create one monitor object per concurrent task set. Monitors are
// thread safe, but they cannot be used to handle multiple sets of multiple
// independent, concurrent tasks. You should use one monitor per asynchronous
// task.
func NewMonitor() *Monitor {
	mon := new(Monitor)
	mon.created = time.Now()
	mon.inboundData = make(chan perfEvent, 1024)
	mon.perfStack = make([]perfEvent, 0)
	mon.aggregatedDataCumulative = make(map[string]time.Duration)
	mon.aggregatedDataSelf = make(map[string]time.Duration)
	go mon.readInboundData()

	return mon
}
