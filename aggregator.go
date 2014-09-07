package sojourner

import (
	"time"
)

var inboundData chan PerfEvent
var perfStack []PerfEvent
var started bool

var aggregatedDataCumulative map[string]time.Duration
var aggregatedDataSelf map[string]time.Duration

// Start() will initialize the performance aggregator and allow you to begin
// collecting information
func Start() {
	inboundData = make(chan PerfEvent, 1024)
	perfStack = make([]PerfEvent, 0)
	aggregatedDataCumulative = make(map[string]time.Duration)
	aggregatedDataSelf = make(map[string]time.Duration)
	go readInboundData()
	started = true
}

func readInboundData() {
	for {
		inbound := <-inboundData

		switch inbound.GetType() {
		case PERF_START:
			perfStack = append(perfStack, inbound)
		case PERF_END:
			lps := len(perfStack)
			name := inbound.GetName()
			lastItem := perfStack[lps-1]
			perfStack = perfStack[:lps-1]

			if lastItem.GetName() != name {
				panic("Unbalanced performance data")
			}

			cumulativeTime := inbound.GetTimestamp().Sub(lastItem.GetTimestamp())
			selfTime := cumulativeTime - lastItem.GetSubtractedTime()

			if _, ok := aggregatedDataCumulative[name]; !ok {
				aggregatedDataCumulative[name] = cumulativeTime
				aggregatedDataSelf[name] = selfTime

			} else {
				aggregatedDataCumulative[name] += cumulativeTime
				aggregatedDataSelf[name] += selfTime
			}

			if lps >= 2 {
				perfStack[lps-2].SubtractSelf(cumulativeTime)
			}
		}

	}
}
