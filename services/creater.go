package services

import (
	"fmt"
	"os"
)

type CreaterService struct {
}

func NewCreaterService() CreaterService {
	return CreaterService{}
}

func (s CreaterService) CreateProject(name, projectDir string) (projectPath, error) {
	projectPath := fmt.Sprintf("%v%v", projectDir, name)

	err := os.Mkdir(projectPath, 0775) // default permission for folders

	return projectPath, err
}
