package domain

import (
	"testing"
	"time"
)

func TestValueObjects(t *testing.T) {
	s, e := NewRobotSerial(" r-7 ")
	if e != nil || s.String() != "R-7" {
		t.Fatalf("serial %v %s", e, s)
	}
	if _, e = NewRobotSerial("bad space"); e == nil {
		t.Fatal("invalid serial accepted")
	}
	c := Cursor{At: time.Unix(2, 3), ID: "m"}
	round, e := ParseCursor(c.Encode())
	if e != nil || round.ID != "m" {
		t.Fatal(e)
	}
	set := CapabilitySet{"vision": 2}
	copy := set.Clone()
	copy["arm"] = 1
	if len(set) != 1 {
		t.Fatal("clone alias")
	}
}
