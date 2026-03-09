package build

import (
	"testing"
)

func TestJobAppendLog(t *testing.T) {
	job := &Job{
		Logs: make(chan string, 10),
	}

	job.AppendLog("line 1")
	job.AppendLog("line 2")

	history := job.LogHistory()
	if len(history) != 2 {
		t.Fatalf("expected 2 log lines, got %d", len(history))
	}
	if history[0] != "line 1" {
		t.Errorf("expected 'line 1', got %q", history[0])
	}
	if history[1] != "line 2" {
		t.Errorf("expected 'line 2', got %q", history[1])
	}
}

func TestJobMarkDone(t *testing.T) {
	job := &Job{
		Status: StatusRunning,
		Logs:   make(chan string, 10),
	}

	job.MarkDone(StatusSuccess, "")
	if job.Status != StatusSuccess {
		t.Errorf("expected status success, got %s", job.Status)
	}
	if job.EndedAt.IsZero() {
		t.Error("expected EndedAt to be set")
	}

	// Second call should be a no-op
	job.MarkDone(StatusFailed, "err")
	if job.Status != StatusSuccess {
		t.Errorf("expected status to remain success, got %s", job.Status)
	}
}

func TestJobMarkDoneWithError(t *testing.T) {
	job := &Job{
		Status: StatusRunning,
		Logs:   make(chan string, 10),
	}

	job.MarkDone(StatusFailed, "build error")
	if job.Status != StatusFailed {
		t.Errorf("expected status failed, got %s", job.Status)
	}
	if job.Error != "build error" {
		t.Errorf("expected error 'build error', got %q", job.Error)
	}
}
