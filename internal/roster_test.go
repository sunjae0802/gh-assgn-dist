package internal

import (
	"os"
	"path/filepath"
	"reflect"
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

func TestLoadRoster_WithHeader(t *testing.T) {
	path := writeRoster(t, "name,email,github\nAlice,alice@example.com,alice\nBob,bob@example.com,bob\n")

	got, err := LoadRoster(path)
	if err != nil {
		t.Fatalf("LoadRoster() error = %v", err)
	}

	want := []Student{
		{Name: "Alice", Email: "alice@example.com", GitHub: "alice"},
		{Name: "Bob", Email: "bob@example.com", GitHub: "bob"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("LoadRoster() = %+v, want %+v", got, want)
	}
}

func TestLoadRoster_NoHeader(t *testing.T) {
	path := writeRoster(t, "Alice,alice@example.com,alice\n")

	got, err := LoadRoster(path)
	if err != nil {
		t.Fatalf("LoadRoster() error = %v", err)
	}

	want := []Student{{Name: "Alice", Email: "alice@example.com", GitHub: "alice"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("LoadRoster() = %+v, want %+v", got, want)
	}
}

func TestLoadRoster_MismatchedFieldCount(t *testing.T) {
	// encoding/csv enforces a consistent field count per row, so a row with
	// fewer fields than the header is a parse error rather than being skipped.
	path := writeRoster(t, "name,email,github\nBadRow,onlytwo\n")
	if _, err := LoadRoster(path); err == nil {
		t.Fatal("expected error for mismatched field count, got nil")
	}
}

func TestLoadRoster_Empty(t *testing.T) {
	path := writeRoster(t, "")
	if _, err := LoadRoster(path); err == nil {
		t.Fatal("expected error for empty roster file, got nil")
	}
}

func TestLoadRoster_MissingFile(t *testing.T) {
	_, err := LoadRoster(filepath.Join(t.TempDir(), "nope.csv"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
