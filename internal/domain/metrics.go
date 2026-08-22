package domain

import (
	"math"
	"sort"
	"time"
)

type MetricSample struct {
	Name   string
	Value  float64
	At     time.Time
	Labels map[string]string
}
type MetricSeries struct {
	Name    string
	Samples []MetricSample
}

func (s MetricSeries) Clone() MetricSeries {
	out := MetricSeries{Name: s.Name, Samples: make([]MetricSample, len(s.Samples))}
	copy(out.Samples, s.Samples)
	for i := range out.Samples {
		out.Samples[i].Labels = NormalizeLabels(out.Samples[i].Labels)
	}
	return out
}
func (s MetricSeries) Sort() MetricSeries {
	out := s.Clone()
	sort.SliceStable(out.Samples, func(i, j int) bool { return out.Samples[i].At.Before(out.Samples[j].At) })
	return out
}
func (s MetricSeries) Average() float64 {
	if len(s.Samples) == 0 {
		return 0
	}
	var total float64
	for _, v := range s.Samples {
		total += v.Value
	}
	return total / float64(len(s.Samples))
}
func (s MetricSeries) P95() float64 {
	if len(s.Samples) == 0 {
		return 0
	}
	v := s.Sort().Samples
	idx := int(math.Ceil(float64(len(v))*0.95)) - 1
	if idx < 0 {
		idx = 0
	}
	return v[idx].Value
}
func (s MetricSeries) Window(start, end time.Time) MetricSeries {
	out := MetricSeries{Name: s.Name}
	for _, v := range s.Samples {
		if !v.At.Before(start) && v.At.Before(end) {
			out.Samples = append(out.Samples, v)
		}
	}
	return out
}
