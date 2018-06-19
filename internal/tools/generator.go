package tools

// NewGenerator starts foreground goroutine which generates sequence of unsigned ints and
// puts them in input channel, also it returnes stop channel which need to be triggered when
// generator need to be stopped
func NewGenerator(input chan<- uint) chan<- bool {
	stop := make(chan bool)

	go func() {
		var current uint = 1
		for {
			select {
			case input <- current:
				current++
			case <-stop:
				close(input)
				return
			}
		}
	}()

	return stop
}
