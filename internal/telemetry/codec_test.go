package telemetry

import (
	"encoding/json"
	"testing"
)

func TestCodec(t *testing.T) {
	b, e := Encode(Envelope{Version: 1, Robot: "r", Sequence: 1, Payload: json.RawMessage(`{"x":1}`)})
	if e != nil {
		t.Fatal(e)
	}
	out, e := Decode(b)
	if e != nil || out.Robot != "r" {
		t.Fatalf("%v %+v", e, out)
	}
}
