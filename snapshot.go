package sojourner

import (
	"sort"
	"time"
)

// A performance snapshot is a collection of all of the data collected by a
// Monitor at any given point in time.
type PerformanceSnapshot struct {
	// All of the metrics from the current snapshot, ordered from least to most
	// substantial.
	Metrics []PerformanceMetric

	created time.Time
}

// A performance metric is a single facet of a snapshot. It represents a single
// type of performance event that took place (initiated by Begin and End).
type PerformanceMetric struct {
	// The name of the represented event
	Name string
	// The amount of time consumed by this type of event and all sub-events.
	Cumulative time.Duration
	// The amount of time consumed by this type of event alone.
	Self time.Duration
}

type metricBag []PerformanceMetric

func (self metricBag) Len() int           { return len(self) }
func (self metricBag) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }
func (self metricBag) Less(i, j int) bool { return self[i].Name < self[j].Name }

// Returns a PerformanceSnapshot of the monitor object. This operation locks
// the monitor so no new events will be processed while the data is collected.
// The processing of the data once it is collected will not block the monitor.
//
// Multiple snapshots can be safely taken in parallel.
func (self *Monitor) Snapshot() PerformanceSnapshot {
	self.aggregateLock.RLock()
	metrics := self.takeSnapshot()
	self.aggregateLock.RUnlock()

	return perfMapToSnapshot(metrics, self.created)
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

func perfMapToSnapshot(metrics map[string]PerformanceMetric, created time.Time) PerformanceSnapshot {
	result := make([]PerformanceMetric, len(metrics))
	for _, metric := range metrics {
		result = append(result, metric)
	}

	sort.Sort(metricBag(result))

	converted := PerformanceSnapshot{result, created}
	return converted
}

// CombineSnapshot combines two or more PerformanceSnapshot objects into a
// single PerformanceSnapshot object. This is done by combining individual
// PerformanceMetric objects with identical nmes.
func CombineSnapshots(snapshots []PerformanceSnapshot) PerformanceSnapshot {
	metrics := make(
		map[string]PerformanceMetric,
		1024,
	)

	for _, snapshot := range snapshots {
		for _, metric := range snapshot.Metrics {
			if existingMetric, ok := metrics[metric.Name]; ok {
				existingMetric.Cumulative += metric.Cumulative
				existingMetric.Self += metric.Self
			} else {
				clone := PerformanceMetric{
					metric.Name,
					metric.Cumulative,
					metric.Self,
				}
				metrics[metric.Name] = clone
			}
		}
	}

	return perfMapToSnapshot(metrics, snapshots[0].created)
}
