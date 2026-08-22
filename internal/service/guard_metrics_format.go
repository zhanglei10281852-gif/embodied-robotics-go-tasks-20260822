package service

import (
	"encoding/json"
	"fmt"
)

func FormatGuardMetrics(m *GuardMetrics) string {
	a, b, c, last := m.Snapshot()
	v := map[string]any{"allowed": a, "denied": b, "expired": c, "last": last}
	raw, e := json.Marshal(v)
	if e != nil {
		return fmt.Sprintf("guard metrics unavailable: %v", e)
	}
	return string(raw)
}
