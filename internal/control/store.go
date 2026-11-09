package control

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

type Store struct {
	path string
}

type Project struct {
	Version int    `json:"version"`
	Name    string `json:"name"`
	Tasks   []Task `json:"tasks"`
}

type Task struct {
	ID                  string      `json:"id"`
	Repo                string      `json:"repo"`
	Title               string      `json:"title,omitempty"`
	Objective           string      `json:"objective,omitempty"`
	Acceptance          []string    `json:"acceptance,omitempty"`
	HandoffRequirements []string    `json:"handoff_requirements,omitempty"`
	Status              string      `json:"status"`
	DependsOn           []string    `json:"depends_on"`
	Owner               string      `json:"owner,omitempty"`
	Result              *Completion `json:"result,omitempty"`
	Resolution          string      `json:"resolution,omitempty"`
}

type Completion struct {
	Commit     string   `json:"commit"`
	Validation []string `json:"validation"`
	Handoff    string   `json:"handoff"`
}

type Decision struct {
	ID            string `json:"id"`
	TaskID        string `json:"task_id"`
	Original      string `json:"original"`
	Problem       string `json:"problem"`
	Choice        string `json:"choice"`
	NeedsApproval bool   `json:"needs_approval"`
}

type Summary struct {
	Counts map[string]int `json:"counts"`
	Ready  []string       `json:"ready"`
}

type TaskContext struct {
	Task         Task   `json:"task"`
	Dependencies []Task `json:"dependencies"`
}

func NewStore(path string) *Store {
	return &Store{path: path}
}

func (s *Store) Load() (*Project, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return nil, err
	}
	var project Project
	if err := json.Unmarshal(data, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

func (s *Store) Summary() (Summary, error) {
	project, err := s.Load()
	if err != nil {
		return Summary{}, err
	}

	statusByID := make(map[string]string, len(project.Tasks))
	for _, task := range project.Tasks {
		statusByID[task.ID] = task.Status
	}

	summary := Summary{Counts: make(map[string]int), Ready: []string{}}
	for _, task := range project.Tasks {
		summary.Counts[task.Status]++
		if task.Status != "pending" {
			continue
		}
		ready := true
		for _, dependencyID := range task.DependsOn {
			status, ok := statusByID[dependencyID]
			if !ok || (status != "completed" && status != "skipped") {
				ready = false
				break
			}
		}
		if ready {
			summary.Ready = append(summary.Ready, task.ID)
		}
	}
	return summary, nil
}

func (s *Store) Context(taskID string) (TaskContext, error) {
	project, err := s.Load()
	if err != nil {
		return TaskContext{}, err
	}
	tasksByID := make(map[string]Task, len(project.Tasks))
	for _, task := range project.Tasks {
		tasksByID[task.ID] = task
	}
	task, ok := tasksByID[taskID]
	if !ok {
		return TaskContext{}, fmt.Errorf("task %q not found", taskID)
	}
	context := TaskContext{Task: task}
	for _, dependencyID := range task.DependsOn {
		dependency, ok := tasksByID[dependencyID]
		if !ok {
			return TaskContext{}, fmt.Errorf("task %q depends on unknown task %q", taskID, dependencyID)
		}
		context.Dependencies = append(context.Dependencies, dependency)
	}
	return context, nil
}

func (s *Store) Claim(taskID, owner string) error {
	return s.withLock(func() error {
		project, err := s.Load()
		if err != nil {
			return err
		}
		statusByID := make(map[string]string, len(project.Tasks))
		for _, task := range project.Tasks {
			statusByID[task.ID] = task.Status
		}
		for i := range project.Tasks {
			if project.Tasks[i].ID != taskID {
				continue
			}
			if project.Tasks[i].Status != "pending" {
				return fmt.Errorf("task %q has status %q and cannot be claimed", taskID, project.Tasks[i].Status)
			}
			for _, dependencyID := range project.Tasks[i].DependsOn {
				status, ok := statusByID[dependencyID]
				if !ok {
					return fmt.Errorf("task %q depends on unknown task %q", taskID, dependencyID)
				}
				if status != "completed" && status != "skipped" {
					return fmt.Errorf("task %q dependency %q has status %q", taskID, dependencyID, status)
				}
			}
			project.Tasks[i].Status = "claimed"
			project.Tasks[i].Owner = owner
			return s.save(project)
		}
		return fmt.Errorf("task %q not found", taskID)
	})
}

func (s *Store) Complete(taskID, owner string, completion Completion) error {
	if completion.Commit == "" {
		return fmt.Errorf("completion commit is required")
	}
	if len(completion.Validation) == 0 {
		return fmt.Errorf("completion validation evidence is required")
	}
	if completion.Handoff == "" {
		return fmt.Errorf("completion handoff is required")
	}
	return s.withLock(func() error {
		project, err := s.Load()
		if err != nil {
			return err
		}
		for i := range project.Tasks {
			if project.Tasks[i].ID != taskID {
				continue
			}
			if project.Tasks[i].Owner != owner {
				return fmt.Errorf("task %q is owned by %q", taskID, project.Tasks[i].Owner)
			}
			project.Tasks[i].Status = "completed"
			project.Tasks[i].Result = &completion
			return s.save(project)
		}
		return fmt.Errorf("task %q not found", taskID)
	})
}

func (s *Store) Skip(taskID, reason string) error {
	if reason == "" {
		return fmt.Errorf("skip reason is required")
	}
	return s.withLock(func() error {
		project, err := s.Load()
		if err != nil {
			return err
		}
		for i := range project.Tasks {
			if project.Tasks[i].ID != taskID {
				continue
			}
			project.Tasks[i].Status = "skipped"
			project.Tasks[i].Resolution = reason
			return s.save(project)
		}
		return fmt.Errorf("task %q not found", taskID)
	})
}

func (s *Store) withLock(action func() error) error {
	lock, err := os.OpenFile(s.path+".lock", os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return err
	}
	defer lock.Close()
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX); err != nil {
		return err
	}
	defer syscall.Flock(int(lock.Fd()), syscall.LOCK_UN) //nolint:errcheck
	return action()
}

func (s *Store) save(project *Project) error {
	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return err
	}
	temporary, err := os.CreateTemp(filepath.Dir(s.path), ".project-*.tmp")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)
	if err := temporary.Chmod(0o600); err != nil {
		temporary.Close()
		return err
	}
	if _, err := temporary.Write(append(data, '\n')); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	return os.Rename(temporaryPath, s.path)
}

func (s *Store) AddDecision(decision Decision) error {
	data, err := json.Marshal(decision)
	if err != nil {
		return err
	}
	path := filepath.Join(filepath.Dir(s.path), "decisions.jsonl")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(append(data, '\n'))
	return err
}
