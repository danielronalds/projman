package services

import "os/exec"

type HealthService struct{}

type RequirementResult struct {
	Name      string
	Installed bool
}

func NewHealthService() HealthService {
	return HealthService{}
}

func (s HealthService) CheckRequirements(programs []string) map[string]bool {
	status := make(map[string]bool, len(programs))

	for _, program := range programs {
		_, err := exec.LookPath(program)
		status[program] = err == nil
	}

	return status
}
