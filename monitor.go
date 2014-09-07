package sojourner

import (
	"sync"
	"time"
)

type Monitor struct {
	inboundData chan PerfEvent
	perfStack   []PerfEvent

	aggregatedDataCumulative map[string]time.Duration
	aggregatedDataSelf       map[string]time.Duration

	aggregateLock sync.RWMutex
}

func NewMonitor() *Monitor {
	mon := new(Monitor)
	mon.inboundData = make(chan PerfEvent, 1024)
	mon.perfStack = make([]PerfEvent, 0)
	mon.aggregatedDataCumulative = make(map[string]time.Duration)
	mon.aggregatedDataSelf = make(map[string]time.Duration)
	go mon.readInboundData()

	return mon
}
