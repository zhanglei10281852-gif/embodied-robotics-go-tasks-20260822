package service

import "time"

const GuardMetricsVersion = "v1"
const GuardMetricAllowed = "allowed"
const GuardMetricDenied = "denied"
const GuardMetricExpired = "expired"
const GuardMetricHealthy = "healthy"
const GuardMetricWindow = "window"
const GuardMetricRate = "rate"

func GuardMetricsHealthy(m *GuardMetrics, since time.Time) bool {
	allowed, denied, expired, last := m.Snapshot()
	if allowed == 0 && denied == 0 && expired == 0 {
		return true
	}
	if last.IsZero() || last.Before(since) {
		return false
	}
	return denied+expired <= allowed*10+1
}

func GuardMetricsRate(m *GuardMetrics) float64 {
	a, d, e, _ := m.Snapshot()
	denom := a + d + e
	if denom == 0 {
		return 0
	}
	return float64(a) / float64(denom)
}

func GuardMetricsWindow(m *GuardMetrics, start, end time.Time) bool {
	_, _, _, last := m.Snapshot()
	return !last.IsZero() && !last.Before(start) && last.Before(end)
}
