package controllers

import (
	"fmt"
	"testing"

	"github.com/danielronalds/projman/internal/services"
)

type mockGroupedLister struct {
	groups    []services.ProjectGroup
	returnErr error
}

func (m mockGroupedLister) ListProjectsByDirectory(filter string) ([]services.ProjectGroup, error) {
	return m.groups, m.returnErr
}

func TestListController_HandleArgs_ServiceError(t *testing.T) {
	mock := mockGroupedLister{returnErr: fmt.Errorf("config broken")}
	controller := NewListController(mock)

	err := controller.HandleArgs([]string{"list"})
	if err == nil {
		t.Fatalf("HandleArgs() error = nil, want error")
	}
}
