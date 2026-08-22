package middleware

import "testing"

func TestBearerRejectsAmbiguousAuthorization(t *testing.T) {
	if _, err := Bearer("Bearer token trailing"); err == nil {
		t.Fatal("ambiguous bearer header accepted")
	}
}
