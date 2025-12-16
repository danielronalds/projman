package controllers

import "testing"

func TestParseRepoName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "httpsWithGit", input: "https://github.com/org/repo.git", want: "repo"},
		{name: "trailingSlash", input: "https://github.com/org/repo.git/", want: "repo"},
		{name: "scpFormat", input: "git@github.com:org/repo.git", want: "repo"},
		{name: "withoutGitExtension", input: "https://github.com/org/repo", want: "repo"},
		{name: "simpleName", input: "repo", want: "repo"},
		{name: "withWhitespace", input: "  git@github.com:org/repo.git  ", want: "repo"},
		{name: "invalid", input: "   ", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRepoName(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("parseRepoName() error = nil, want error")
				}
				return
			}

			if err != nil {
				t.Fatalf("parseRepoName() error = %v, want nil", err)
			}

			if got != tt.want {
				t.Fatalf("parseRepoName() = %q, want %q", got, tt.want)
			}
		})
	}
}
