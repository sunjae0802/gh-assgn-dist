package cmd

import (
	"testing"

	"github.com/sunjae0802/gh-assgn-dist/internal"
)

func TestSplitRepo(t *testing.T) {
	tests := []struct {
		repo      string
		wantOwner string
		wantName  string
		wantErr   bool
	}{
		{repo: "witcomp1000/a1-template", wantOwner: "witcomp1000", wantName: "a1-template"},
		{repo: "owner/name/extra", wantOwner: "owner", wantName: "name/extra"},
		{repo: "no-slash", wantErr: true},
		{repo: "", wantErr: true},
	}

	for _, tt := range tests {
		owner, name, err := splitRepo(tt.repo)
		if tt.wantErr {
			if err == nil {
				t.Errorf("splitRepo(%q): expected error, got nil", tt.repo)
			}
			continue
		}
		if err != nil {
			t.Errorf("splitRepo(%q): unexpected error: %v", tt.repo, err)
			continue
		}
		if owner != tt.wantOwner || name != tt.wantName {
			t.Errorf("splitRepo(%q) = (%q, %q), want (%q, %q)", tt.repo, owner, name, tt.wantOwner, tt.wantName)
		}
	}
}

func TestFilterStudents(t *testing.T) {
	roster := []internal.Student{
		{Name: "Alice", GitHub: "alice"},
		{Name: "Bob", GitHub: "Bob"},
		{Name: "Carol", GitHub: "carol"},
	}

	tests := []struct {
		name          string
		only          string
		wantGitHub    []string
		wantUnmatched []string
	}{
		{name: "empty selects everyone", only: "", wantGitHub: []string{"alice", "Bob", "carol"}},
		{name: "whitespace selects everyone", only: "  ", wantGitHub: []string{"alice", "Bob", "carol"}},
		{name: "single", only: "alice", wantGitHub: []string{"alice"}},
		{name: "multiple", only: "alice,carol", wantGitHub: []string{"alice", "carol"}},
		{name: "spaces around names", only: " alice , carol ", wantGitHub: []string{"alice", "carol"}},
		{name: "case-insensitive", only: "ALICE,bob", wantGitHub: []string{"alice", "Bob"}},
		{name: "empty entries ignored", only: "alice,,carol,", wantGitHub: []string{"alice", "carol"}},
		{name: "duplicates select once", only: "alice,alice", wantGitHub: []string{"alice"}},
		{name: "unknown reported", only: "alice,dave", wantGitHub: []string{"alice"}, wantUnmatched: []string{"dave"}},
		{name: "all unknown", only: "dave", wantUnmatched: []string{"dave"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, unmatched := filterStudents(roster, tt.only)

			var gotGitHub []string
			for _, s := range got {
				gotGitHub = append(gotGitHub, s.GitHub)
			}
			if !equalStrings(gotGitHub, tt.wantGitHub) {
				t.Errorf("selected = %v, want %v", gotGitHub, tt.wantGitHub)
			}
			if !equalStrings(unmatched, tt.wantUnmatched) {
				t.Errorf("unmatched = %v, want %v", unmatched, tt.wantUnmatched)
			}
		})
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
