package cmd

import "testing"

func TestPluralStudents(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{n: 0, want: "0 students"},
		{n: 1, want: "1 student"},
		{n: 2, want: "2 students"},
		{n: 42, want: "42 students"},
	}

	for _, tt := range tests {
		if got := pluralStudents(tt.n); got != tt.want {
			t.Errorf("pluralStudents(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}
