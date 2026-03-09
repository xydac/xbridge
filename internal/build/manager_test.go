package build

import (
	"fmt"
	"testing"

	"github.com/xydac/xbridge/internal/config"
)

func TestManagerGetJobNotFound(t *testing.T) {
	mgr := NewManager(t.TempDir(), config.DefaultConfig())

	_, ok := mgr.GetJob("nonexistent")
	if ok {
		t.Error("expected job not found")
	}
}

func TestManagerPruneOldJobs(t *testing.T) {
	mgr := NewManager(t.TempDir(), config.DefaultConfig())

	// Add 7 jobs manually with unique IDs
	for i := 0; i < 7; i++ {
		id := fmt.Sprintf("build-test-%d", i)
		mgr.jobs[id] = &Job{ID: id, Logs: make(chan string, 1)}
		mgr.jobOrder = append(mgr.jobOrder, id)
	}

	mgr.pruneOldJobs()

	if len(mgr.jobs) != maxKeptJobs {
		t.Errorf("expected %d jobs after prune, got %d", maxKeptJobs, len(mgr.jobs))
	}
	if len(mgr.jobOrder) != maxKeptJobs {
		t.Errorf("expected %d job order entries, got %d", maxKeptJobs, len(mgr.jobOrder))
	}
}
