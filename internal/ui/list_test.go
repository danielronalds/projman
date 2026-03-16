package ui

import (
	"strings"
	"testing"

	"github.com/danielronalds/projman/internal/services"
)

func TestRenderProjectList(t *testing.T) {
	tests := []struct {
		name          string
		groups        []services.ProjectGroup
		filtered      bool
		wantContains  []string
		wantExcludes  []string
		wantCountLine string
	}{
		{
			name: "multiple groups rendered with bold headers and bullets",
			groups: []services.ProjectGroup{
				{Directory: "Personal/", Projects: []string{"alpha", "bravo"}},
				{Directory: "Work/", Projects: []string{"charlie"}},
			},
			wantContains: []string{
				"\033[1mPersonal/\033[0m",
				" \u2022 alpha",
				" \u2022 bravo",
				"\033[1mWork/\033[0m",
				" \u2022 charlie",
			},
			wantCountLine: "3 projects total on your system",
		},
		{
			name: "single project uses singular label",
			groups: []services.ProjectGroup{
				{Directory: "Work/", Projects: []string{"only-one"}},
			},
			wantContains: []string{
				" \u2022 only-one",
			},
			wantCountLine: "1 project total on your system",
		},
		{
			name:          "no groups shows zero count",
			groups:        []services.ProjectGroup{},
			wantCountLine: "0 projects total on your system",
		},
		{
			name: "filtered footer says matched instead of on your system",
			groups: []services.ProjectGroup{
				{Directory: "Work/", Projects: []string{"alpha", "bravo"}},
			},
			filtered:      true,
			wantCountLine: "2 projects matched",
		},
		{
			name:          "filtered with no matches",
			groups:        []services.ProjectGroup{},
			filtered:      true,
			wantCountLine: "0 projects matched",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := RenderProjectList(tt.groups, tt.filtered)

			for _, s := range tt.wantContains {
				if !strings.Contains(output, s) {
					t.Errorf("output missing %q\ngot:\n%s", s, output)
				}
			}

			for _, s := range tt.wantExcludes {
				if strings.Contains(output, s) {
					t.Errorf("output should not contain %q\ngot:\n%s", s, output)
				}
			}

			if tt.wantCountLine != "" && !strings.Contains(output, tt.wantCountLine) {
				t.Errorf("output missing count line %q\ngot:\n%s", tt.wantCountLine, output)
			}
		})
	}
}
