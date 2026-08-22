package domain

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTimestampAndResult(t *testing.T) {
	v := Timestamp{Time: time.Unix(2, 0).UTC()}
	b, e := json.Marshal(v)
	if e != nil {
		t.Fatal(e)
	}
	var out Timestamp
	if e = json.Unmarshal(b, &out); e != nil || !out.Equal(v.Time) {
		t.Fatalf("%v %v", e, out)
	}
	if e = (CommandResult{Accepted: true}).Validate(); e == nil {
		t.Fatal("missing queue")
	}
}
