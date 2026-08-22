package worker

import "time"

type Backoff struct {
	Base time.Duration
	Max  time.Duration
}

func (b Backoff) Delay(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	d := b.Base
	for i := 0; i < attempt; i++ {
		d *= 2
		if d >= b.Max {
			return b.Max
		}
	}
	if d > b.Max {
		return b.Max
	}
	return d
}
