package build

import (
	"sync"
	"time"
)

// Status represents the state of a build job.
type Status string

const (
	StatusQueued  Status = "queued"
	StatusRunning Status = "running"
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
)

// Job represents a single build job.
type Job struct {
	ID        string    `json:"id"`
	Status    Status    `json:"status"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at,omitempty"`
	Error     string    `json:"error,omitempty"`
	Profile   string    `json:"profile,omitempty"`

	// Logs is a channel for streaming build output.
	// Consumers should range over this channel.
	Logs chan string `json:"-"`

	// logHistory stores all log lines for replay.
	mu         sync.Mutex
	logHistory []string
	done       bool
}

// AppendLog adds a log line to history and sends it to the channel.
func (j *Job) AppendLog(line string) {
	j.mu.Lock()
	j.logHistory = append(j.logHistory, line)
	j.mu.Unlock()

	// Non-blocking send
	select {
	case j.Logs <- line:
	default:
	}
}

// LogHistory returns a copy of all log lines.
func (j *Job) LogHistory() []string {
	j.mu.Lock()
	defer j.mu.Unlock()
	out := make([]string, len(j.logHistory))
	copy(out, j.logHistory)
	return out
}

// MarkDone marks the job as done and closes the log channel.
func (j *Job) MarkDone(status Status, errMsg string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.done {
		return
	}
	j.done = true
	j.Status = status
	j.Error = errMsg
	j.EndedAt = time.Now()
	close(j.Logs)
}
