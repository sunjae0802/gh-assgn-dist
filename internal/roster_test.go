package internal

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func writeRoster(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "roster.csv")
	if err := os.WriteFile(path, []byte(contents), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func loadRoster(t *testing.T, contents string) ([]Student, []string) {
	t.Helper()
	students, warnings, err := LoadRoster(writeRoster(t, contents))
	if err != nil {
		t.Fatalf("LoadRoster() error = %v", err)
	}
	return students, warnings
}

func TestLoadRoster_WithHeader(t *testing.T) {
	got, warnings := loadRoster(t, "name,email,github\nAlice,alice@example.com,alice\nBob,bob@example.com,bob\n")

	want := []Student{
		{Name: "Alice", Email: "alice@example.com", GitHub: "alice"},
		{Name: "Bob", Email: "bob@example.com", GitHub: "bob"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("LoadRoster() = %+v, want %+v", got, want)
	}
	if len(warnings) != 0 {
		t.Errorf("warnings = %v, want none", warnings)
	}
}

func TestLoadRoster_GitHubHeaderCapitalization(t *testing.T) {
	for _, header := range []string{"github", "GitHub", "Github", "GITHUB", " GitHub "} {
		t.Run(header, func(t *testing.T) {
			got, _ := loadRoster(t, "name,email,"+header+"\nAlice,alice@example.com,alice\n")

			want := []Student{{Name: "Alice", Email: "alice@example.com", GitHub: "alice"}}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("LoadRoster() = %+v, want %+v", got, want)
			}
		})
	}
}

func TestLoadRoster_ColumnOrderIndependent(t *testing.T) {
	got, _ := loadRoster(t, "GitHub,Name,Email\nalice,Alice,alice@example.com\n")

	want := []Student{{Name: "Alice", Email: "alice@example.com", GitHub: "alice"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("LoadRoster() = %+v, want %+v", got, want)
	}
}

func TestLoadRoster_GitHubColumnAlone(t *testing.T) {
	got, warnings := loadRoster(t, "github\nalice\n")

	want := []Student{{GitHub: "alice"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("LoadRoster() = %+v, want %+v", got, want)
	}
	if len(warnings) != 0 {
		t.Errorf("warnings = %v, want none", warnings)
	}
}

func TestLoadRoster_MissingGitHubColumn(t *testing.T) {
	for _, contents := range []string{
		"name,email\nAlice,alice@example.com\n",
		"Alice,alice@example.com,alice\n", // headerless: first row read as the header
		"name,email,gituhb\nAlice,alice@example.com,alice\n",
	} {
		_, _, err := LoadRoster(writeRoster(t, contents))
		if err == nil {
			t.Fatalf("expected error for roster without a github column: %q", contents)
		}
		if !strings.Contains(err.Error(), "github") {
			t.Errorf("error = %v, want it to mention the missing github column", err)
		}
	}
}

func TestLoadRoster_SkipsBlankRows(t *testing.T) {
	got, warnings := loadRoster(t, "name,email,github\nAlice,alice@example.com,alice\n\n,,\n   ,,\nBob,bob@example.com,bob\n")

	want := []Student{
		{Name: "Alice", Email: "alice@example.com", GitHub: "alice"},
		{Name: "Bob", Email: "bob@example.com", GitHub: "bob"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("LoadRoster() = %+v, want %+v", got, want)
	}
	if len(warnings) != 0 {
		t.Errorf("warnings = %v, want none for blank rows", warnings)
	}
}

func TestLoadRoster_WarnsOnMissingUsername(t *testing.T) {
	tests := []struct {
		name     string
		row      string
		wantWarn string
	}{
		{name: "name and email", row: "Carol,carol@example.com,", wantWarn: "Carol <carol@example.com>"},
		{name: "name only", row: "Carol,,", wantWarn: "Carol"},
		{name: "email only", row: ",carol@example.com,", wantWarn: "carol@example.com"},
		{name: "short row", row: "Carol,carol@example.com", wantWarn: "Carol <carol@example.com>"},
		{name: "whitespace username", row: "Carol,carol@example.com,   ", wantWarn: "Carol <carol@example.com>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Blank line before the bad row: the warning must cite the
			// physical line, not the row's ordinal.
			got, warnings := loadRoster(t, "name,email,github\nAlice,alice@example.com,alice\n\n"+tt.row+"\n")

			want := []Student{{Name: "Alice", Email: "alice@example.com", GitHub: "alice"}}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("LoadRoster() = %+v, want %+v (row must be skipped)", got, want)
			}
			if len(warnings) != 1 {
				t.Fatalf("warnings = %v, want exactly one", warnings)
			}
			if !strings.Contains(warnings[0], tt.wantWarn) {
				t.Errorf("warning = %q, want it to name %q", warnings[0], tt.wantWarn)
			}
			if !strings.Contains(warnings[0], "line 4") {
				t.Errorf("warning = %q, want it to cite line 4", warnings[0])
			}
		})
	}
}

func TestLoadRoster_StripsBOM(t *testing.T) {
	got, _ := loadRoster(t, "\ufeffname,email,github\nAlice,alice@example.com,alice\n")

	want := []Student{{Name: "Alice", Email: "alice@example.com", GitHub: "alice"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("LoadRoster() = %+v, want %+v", got, want)
	}
}

func TestLoadRoster_Empty(t *testing.T) {
	if _, _, err := LoadRoster(writeRoster(t, "")); err == nil {
		t.Fatal("expected error for empty roster file, got nil")
	}
}

func TestLoadRoster_HeaderOnly(t *testing.T) {
	got, warnings := loadRoster(t, "name,email,github\n")
	if len(got) != 0 {
		t.Errorf("LoadRoster() = %+v, want no students", got)
	}
	if len(warnings) != 0 {
		t.Errorf("warnings = %v, want none", warnings)
	}
}

func TestLoadRoster_MissingFile(t *testing.T) {
	_, _, err := LoadRoster(filepath.Join(t.TempDir(), "nope.csv"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
