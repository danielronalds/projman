package services

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type mockProjectsConfig struct {
	projectDirs []string
}

func (m mockProjectsConfig) ProjectDirs() []string {
	return m.projectDirs
}

func TestListProjectsByDirectory(t *testing.T) {
	tests := []struct {
		name           string
		setupDirs      [][]string
		filter         string
		expectedGroups []ProjectGroup
	}{
		{
			name: "multiple directories with projects",
			setupDirs: [][]string{
				{"alpha", "charlie", "bravo"},
				{"zulu", "echo"},
			},
			expectedGroups: []ProjectGroup{
				{Directory: "dir0/", Projects: []string{"alpha", "bravo", "charlie"}},
				{Directory: "dir1/", Projects: []string{"echo", "zulu"}},
			},
		},
		{
			name: "single directory",
			setupDirs: [][]string{
				{"foo", "bar"},
			},
			expectedGroups: []ProjectGroup{
				{Directory: "dir0/", Projects: []string{"bar", "foo"}},
			},
		},
		{
			name:      "empty directory",
			setupDirs: [][]string{{}},
			expectedGroups: []ProjectGroup{
				{Directory: "dir0/", Projects: []string{}},
			},
		},
		{
			name: "display name trims home directory prefix",
			setupDirs: [][]string{
				{"project-one"},
			},
			expectedGroups: []ProjectGroup{
				{Directory: "dir0/", Projects: []string{"project-one"}},
			},
		},
		{
			name: "filter matches one directory case insensitive",
			setupDirs: [][]string{
				{"alpha"},
				{"bravo"},
			},
			filter: "dir1",
			expectedGroups: []ProjectGroup{
				{Directory: "dir1/", Projects: []string{"bravo"}},
			},
		},
		{
			name: "filter matches nothing returns empty",
			setupDirs: [][]string{
				{"alpha"},
				{"bravo"},
			},
			filter:         "nonexistent",
			expectedGroups: nil,
		},
		{
			name: "filter partial match",
			setupDirs: [][]string{
				{"alpha"},
				{"bravo"},
			},
			filter: "dir0",
			expectedGroups: []ProjectGroup{
				{Directory: "dir0/", Projects: []string{"alpha"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir := t.TempDir()
			t.Setenv("HOME", homeDir)

			var dirs []string
			for i, projects := range tt.setupDirs {
				dirName := filepath.Join(homeDir, fmt.Sprintf("dir%d", i))
				err := os.MkdirAll(dirName, 0755)
				if err != nil {
					t.Fatalf("failed to create dir: %v", err)
				}
				for _, p := range projects {
					err := os.MkdirAll(filepath.Join(dirName, p), 0755)
					if err != nil {
						t.Fatalf("failed to create project dir: %v", err)
					}
				}
				dirs = append(dirs, dirName+"/")
			}

			config := mockProjectsConfig{projectDirs: dirs}
			service := NewProjectsService(config)

			groups, err := service.ListProjectsByDirectory(tt.filter)
			if err != nil {
				t.Fatalf("ListProjectsByDirectory() error = %v, want nil", err)
			}

			if len(groups) != len(tt.expectedGroups) {
				t.Fatalf("got %d groups, want %d", len(groups), len(tt.expectedGroups))
			}

			for i, expected := range tt.expectedGroups {
				got := groups[i]
				if got.Directory != expected.Directory {
					t.Errorf("group[%d].Directory = %q, want %q", i, got.Directory, expected.Directory)
				}

				if len(got.Projects) != len(expected.Projects) {
					t.Fatalf("group[%d] got %d projects, want %d", i, len(got.Projects), len(expected.Projects))
				}

				for j, name := range expected.Projects {
					if got.Projects[j] != name {
						t.Errorf("group[%d].Projects[%d] = %q, want %q", i, j, got.Projects[j], name)
					}
				}
			}
		})
	}
}
