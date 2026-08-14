package cmd

import "testing"

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
