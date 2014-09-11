package sojourner

import "bytes"
import "fmt"
import "time"

func (self PerformanceSnapshot) Report() string {
	var buf bytes.Buffer
	buf.WriteString("==============================\n")
	buf.WriteString(fmt.Sprintf(
		"%s for %s\n",
		self.created.String(),
		time.Now().Sub(self.created).String(),
	))
	buf.WriteString("==============================\n")
	for _, metric := range self.Metrics {
		buf.WriteString(fmt.Sprintf(
			"%32s%16s%16s\n",
			metric.Name,
			metric.Self,
			metric.Cumulative,
		))
	}
	buf.WriteString("==============================\n")
	return buf.String()
}
