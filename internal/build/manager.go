package build

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/xydac/xbridge/internal/config"
	"github.com/xydac/xbridge/internal/engine"
)

const maxKeptJobs = 5

// Manager handles build job lifecycle.
type Manager struct {
	mu       sync.Mutex
	jobs     map[string]*Job
	jobOrder []string // ordered list of job IDs for pruning
	workDir  string
	cfg      *config.Config
}

// NewManager creates a new build manager.
func NewManager(workDir string, cfg *config.Config) *Manager {
	return &Manager{
		jobs:    make(map[string]*Job),
		workDir: workDir,
		cfg:     cfg,
	}
}

// StartBuild kicks off a new build and returns the job.
func (m *Manager) StartBuild(ctx context.Context, profileName string) (*Job, error) {
	resolved, err := m.cfg.ResolveProfile(profileName)
	if err != nil {
		return nil, err
	}

	job := &Job{
		ID:        generateID(),
		Status:    StatusRunning,
		StartedAt: time.Now(),
		Profile:   profileName,
		Logs:      make(chan string, 100),
	}

	m.mu.Lock()
	m.jobs[job.ID] = job
	m.jobOrder = append(m.jobOrder, job.ID)
	m.pruneOldJobs()
	m.mu.Unlock()

	// Run pre-build hook
	if resolved.PreBuild != "" {
		hookCmd := exec.CommandContext(ctx, "sh", "-c", resolved.PreBuild)
		hookCmd.Dir = m.workDir
		if out, err := hookCmd.CombinedOutput(); err != nil {
			job.AppendLog(fmt.Sprintf("pre_build hook failed: %s\n%s", err, string(out)))
			job.MarkDone(StatusFailed, "pre_build hook failed")
			return job, nil
		}
	}

	// Build destination
	destination := "platform=iOS Simulator,name=iPhone 16"
	if resolved.Device != "" {
		destination = fmt.Sprintf("platform=iOS Simulator,name=%s", resolved.Device)
		if resolved.Runtime != "" {
			destination += fmt.Sprintf(",OS=%s", resolved.Runtime)
		}
	}

	buildCfg := &engine.XcodeBuildConfig{
		Project:       resolved.Project,
		Scheme:        resolved.Scheme,
		Configuration: resolved.Configuration,
		Destination:   destination,
		DerivedData:   m.workDir + "/DerivedData",
		ExtraArgs:     resolved.BuildArgs,
		Env:           resolved.Env,
	}

	go m.runBuild(ctx, job, buildCfg)

	return job, nil
}

func (m *Manager) runBuild(ctx context.Context, job *Job, cfg *engine.XcodeBuildConfig) {
	cmd := engine.XcodeBuildCommand(ctx, cfg, "build")
	cmd.Dir = m.workDir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		job.MarkDone(StatusFailed, err.Error())
		return
	}
	cmd.Stderr = cmd.Stdout // merge stderr into stdout

	if err := cmd.Start(); err != nil {
		job.MarkDone(StatusFailed, err.Error())
		return
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		job.AppendLog(scanner.Text())
	}

	if err := cmd.Wait(); err != nil {
		job.MarkDone(StatusFailed, err.Error())
		return
	}

	job.MarkDone(StatusSuccess, "")
}

// GetJob returns a job by ID.
func (m *Manager) GetJob(id string) (*Job, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	j, ok := m.jobs[id]
	return j, ok
}

// CleanBuild runs xcodebuild clean.
func (m *Manager) CleanBuild(ctx context.Context) error {
	resolved, err := m.cfg.ResolveProfile("")
	if err != nil {
		return err
	}

	cfg := &engine.XcodeBuildConfig{
		Project:       resolved.Project,
		Scheme:        resolved.Scheme,
		Configuration: resolved.Configuration,
		DerivedData:   m.workDir + "/DerivedData",
	}
	return engine.CleanBuild(ctx, cfg)
}

func (m *Manager) pruneOldJobs() {
	for len(m.jobOrder) > maxKeptJobs {
		oldID := m.jobOrder[0]
		m.jobOrder = m.jobOrder[1:]
		delete(m.jobs, oldID)
	}
}

func generateID() string {
	return fmt.Sprintf("build-%d", time.Now().UnixNano())
}
