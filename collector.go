package sojourner

func assertStarted() {
	if !started {
		panic("Written performance data before Start()")
	}
}

// Begin indicates that you want to start timing some particular performance
// metric of name `name`.
func Begin(name string) {
	assertStarted()

	inboundData <- newPerfEvent(PERF_START, name)
}

// End indicates that the performance metric that you have started timing is
// complete and the timing information about it should be stored.
func End(name string) {
	assertStarted()

	inboundData <- newPerfEvent(PERF_END, name)
}
