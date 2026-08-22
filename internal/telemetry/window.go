package telemetry

import "time"

func CurrentWindow(now time.Time, size time.Duration) Window {
	end := now.UTC()
	return Window{Start: end.Add(-size), End: end}
}
func Split(w Window, n int) []Window {
	if n < 1 {
		return nil
	}
	step := w.End.Sub(w.Start) / time.Duration(n)
	out := make([]Window, n)
	for i := range out {
		out[i] = Window{Start: w.Start.Add(step * time.Duration(i)), End: w.Start.Add(step * time.Duration(i+1))}
	}
	return out
}
