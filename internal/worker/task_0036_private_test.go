package worker

import "testing"

func TestCheckpointSnapshotDoesNotPolluteInternalState(t *testing.T) {
	c := NewCheckpoint()
	if err := c.Advance("robot-1", 7); err != nil {
		t.Fatal(err)
	}
	s := c.Snapshot()
	s["robot-1"] = 99
	s["robot-2"] = 3
	if got, _ := c.Read("robot-1"); got != 7 {
		t.Fatalf("snapshot mutated checkpoint: %d", got)
	}
	if _, ok := c.Read("robot-2"); ok {
		t.Fatal("snapshot inserted robot into checkpoint")
	}
}
