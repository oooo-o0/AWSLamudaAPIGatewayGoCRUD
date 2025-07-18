package job_handler

import (
	"context"
	"errors"
	"fmt"

	"kaigo-insurance-system/batch/interface"
	"kaigo-insurance-system/batch/utils"
)

type JobHandler struct {
	jobs map[string]interface.Job
}

func New() *JobHandler {
	return &JobHandler{
		jobs: make(map[string]interface.Job),
	}
}

// Register adds a job with a key (e.g., "daily_billing")
func (h *JobHandler) Register(name string, job interface.Job) {
	h.jobs[name] = job
}

// Execute runs the job by name
func (h *JobHandler) Execute(ctx context.Context, name string) error {
	job, ok := h.jobs[name]
	if !ok {
		return fmt.Errorf("job '%s' not registered", name)
	}

	utils.Logger.Infof("Executing job: %s", name)
	if err := job.Run(ctx); err != nil {
		return fmt.Errorf("job '%s' failed: %w", name, err)
	}
	return nil
}
