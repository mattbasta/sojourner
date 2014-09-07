package sojourner

import (
	"sort"
	"time"
)

type PerformanceSnapshot struct {
	Metrics []PerformanceMetric
}
type PerformanceMetric struct {
	Name       string
	Cumulative time.Duration
	Self       time.Duration
}

type metricBag []PerformanceMetric

func (self metricBag) Len() int           { return len(self) }
func (self metricBag) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }
func (self metricBag) Less(i, j int) bool { return self[i].Name < self[j].Name }

func (self *Monitor) Snapshot() PerformanceSnapshot {
	self.aggregateLock.RLock()
	metrics := self.takeSnapshot()
	self.aggregateLock.RUnlock()
	return convertMetricsToSnapshot(metrics)
}

func (self *Monitor) takeSnapshot() map[string]PerformanceMetric {
	metrics := make(
		map[string]PerformanceMetric,
		len(self.aggregatedDataCumulative)+len(self.perfStack),
	)

	// Process accumulated data
	for name, cumulativeVal := range self.aggregatedDataCumulative {
		metrics[name] = PerformanceMetric{
			name,
			cumulativeVal,
			self.aggregatedDataSelf[name],
		}
	}

	// Process accumulating data
	now := time.Now()
	var accumulation time.Duration
	for i := len(self.perfStack) - 1; i >= 0; i-- {
		stackItem := self.perfStack[i]
		name := stackItem.GetName()

		cumulative := now.Sub(stackItem.GetTimestamp())
		self := cumulative - accumulation - stackItem.GetSubtractedTime()

		if metric, ok := metrics[name]; ok {
			metric.Cumulative += cumulative
			metric.Self += self
		} else {
			metrics[name] = PerformanceMetric{
				name,
				cumulative,
				self,
			}
			accumulation += cumulative
		}
	}

	return metrics
}

func convertMetricsToSnapshot(metrics map[string]PerformanceMetric) PerformanceSnapshot {
	result := make([]PerformanceMetric, len(metrics))
	for _, metric := range metrics {
		result = append(result, metric)
	}

	sort.Sort(metricBag(result))
	return PerformanceSnapshot{result}
}
