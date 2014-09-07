package sojourner

func (self *Monitor) readInboundData() {
	for {
		inbound := <-self.inboundData

		self.aggregateLock.Lock()

		switch inbound.GetType() {
		case PERF_START:
			self.perfStack = append(self.perfStack, inbound)
		case PERF_END:
			lps := len(self.perfStack)
			name := inbound.GetName()

			lastItem := self.perfStack[lps-1]

			self.perfStack = self.perfStack[:lps-1]

			if lastItem.GetName() != name {
				panic("Unbalanced performance data")
			}

			cumulativeTime := inbound.GetTimestamp().Sub(lastItem.GetTimestamp())
			selfTime := cumulativeTime - lastItem.GetSubtractedTime()

			if _, ok := self.aggregatedDataCumulative[name]; !ok {
				self.aggregatedDataCumulative[name] = cumulativeTime
				self.aggregatedDataSelf[name] = selfTime

			} else {
				self.aggregatedDataCumulative[name] += cumulativeTime
				self.aggregatedDataSelf[name] += selfTime
			}

			if lps >= 2 {
				self.perfStack[lps-2].SubtractSelf(cumulativeTime)
			}
		}
		self.aggregateLock.Unlock()

	}
}
