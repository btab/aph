package project

import "time"

// coalesces all notifications that are spaced < maxWait apart
// triggers callback when the last such notification is received
func debouncedCallbackLoop(notifications chan bool, maxWait time.Duration, f func()) {
	var timer *time.Timer

	for {
		if timer == nil {
			select {
			case <-notifications:
				timer = time.NewTimer(maxWait)
			}
		} else {
			select {
			case <-notifications:
				timer = time.NewTimer(maxWait)
			case <-timer.C:
				timer = nil
				f()
			}
		}
	}
}
