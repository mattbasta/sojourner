package sojourner

// Begin indicates that you want to start timing some particular performance
// metric of name `name`.
//
// The aggregator must have been initialized with Start() before this method
// can be called.
//
// `Begin()` and End() must be called on the same goroutine for any given
// operation.
func (self *Monitor) Begin(name string) {
	self.inboundData <- newPerfEvent(PERF_START, name)
}

// End indicates that the performance metric that you have started timing is
// complete and the timing information about it should be stored.
//
// Performance metrics may not be stopped on a different goroutine than the one
// on which they were started from.
func (self *Monitor) End(name string) {
	self.inboundData <- newPerfEvent(PERF_END, name)
}
