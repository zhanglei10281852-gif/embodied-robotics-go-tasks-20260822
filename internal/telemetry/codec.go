package telemetry

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
)

type Envelope struct {
	Version  int             `json:"version"`
	Robot    string          `json:"robot"`
	Sequence int64           `json:"sequence"`
	Payload  json.RawMessage `json:"payload"`
}

func Encode(e Envelope) ([]byte, error) {
	if e.Version < 1 || e.Robot == "" || e.Sequence < 1 {
		return nil, fmt.Errorf("invalid envelope")
	}
	var buf bytes.Buffer
	z := gzip.NewWriter(&buf)
	if err := json.NewEncoder(z).Encode(e); err != nil {
		return nil, err
	}
	if err := z.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func Decode(data []byte) (Envelope, error) {
	z, e := gzip.NewReader(bytes.NewReader(data))
	if e != nil {
		return Envelope{}, e
	}
	defer z.Close()
	var out Envelope
	if e = json.NewDecoder(z).Decode(&out); e != nil {
		return out, e
	}
	return out, nil
}
