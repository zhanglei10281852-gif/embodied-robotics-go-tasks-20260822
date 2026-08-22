package audit

import (
	"testing"
	"time"
)

func TestRetentionWindow(t *testing.T) {
	r := Retention{Keep: time.Hour}
	if !r.Window(time.Now()).Before(time.Now()) {
		t.Fatal("window")
	}
}
