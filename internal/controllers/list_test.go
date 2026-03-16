package controllers

import (
	"fmt"
	"testing"

	"github.com/danielronalds/projman/internal/services"
)

type mockGroupedLister struct {
	groups       []services.ProjectGroup
	returnErr    error
	calledFilter string
}

func (m *mockGroupedLister) ListProjectsByDirectory(filter string) ([]services.ProjectGroup, error) {
	m.calledFilter = filter
	return m.groups, m.returnErr
}

func TestListController_HandleArgs_ServiceError(t *testing.T) {
	mock := &mockGroupedLister{returnErr: fmt.Errorf("config broken")}
	controller := NewListController(mock)

	err := controller.HandleArgs([]string{"list"})
	if err == nil {
		t.Fatalf("HandleArgs() error = nil, want error")
	}
}

func TestListController_HandleArgs(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedFilter string
	}{
		{
			name:           "no filter argument",
			args:           []string{"list"},
			expectedFilter: "",
		},
		{
			name:           "with filter argument",
			args:           []string{"list", "work"},
			expectedFilter: "work",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockGroupedLister{}
			controller := NewListController(mock)

			err := controller.HandleArgs(tt.args)
			if err != nil {
				t.Fatalf("HandleArgs() error = %v, want nil", err)
			}

			if mock.calledFilter != tt.expectedFilter {
				t.Errorf("filter = %q, want %q", mock.calledFilter, tt.expectedFilter)
			}
		})
	}
}
