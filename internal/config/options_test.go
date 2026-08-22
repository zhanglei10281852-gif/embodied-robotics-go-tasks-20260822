package config

import (
	"testing"
	"time"
)

func TestOptions(t *testing.T) {
	o := Load().Options()
	if e := o.Validate(); e != nil {
		t.Fatal(e)
	}
	if !o.HasOrigin("http://localhost/") {
		t.Fatal("origin")
	}
	if o.ShutdownTimeout != 10*time.Second {
		t.Fatal(o.ShutdownTimeout)
	}
}
