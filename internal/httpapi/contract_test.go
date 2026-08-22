package httpapi

import (
	"net/http/httptest"
	"testing"
	"time"
)

func TestQueryParsing(t *testing.T) {
	r := httptest.NewRequest("GET", "/x?limit=999&before=2025-01-01T00:00:00Z", nil)
	if got := decodeLimit(r, 50, 500); got != 500 {
		t.Fatal(got)
	}
	if got := parseTimeQuery(r, "before", time.Time{}); got.IsZero() {
		t.Fatal("time missing")
	}
}
