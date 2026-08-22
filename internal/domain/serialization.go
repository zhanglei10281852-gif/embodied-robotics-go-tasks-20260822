package domain

import (
	"encoding/json"
	"errors"
	"time"
)

type APIError struct {
	Code, Message, RequestID string
	Retryable                bool
}

func (e APIError) Error() string { return e.Code + ": " + e.Message }

type Timestamp struct{ time.Time }

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.UTC().Format(time.RFC3339Nano))
}
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	var v string
	if e := json.Unmarshal(b, &v); e != nil {
		return e
	}
	parsed, e := time.Parse(time.RFC3339Nano, v)
	if e != nil {
		return e
	}
	t.Time = parsed
	return nil
}

type CommandResult struct {
	Accepted bool
	QueueID  string
	Error    *APIError
	At       time.Time
}

func (r CommandResult) Validate() error {
	if r.Accepted && r.QueueID == "" {
		return errors.New("accepted result missing queue id")
	}
	if !r.Accepted && r.Error == nil {
		return errors.New("rejected result missing error")
	}
	return nil
}
