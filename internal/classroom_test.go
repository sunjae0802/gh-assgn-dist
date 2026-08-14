package internal

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestRepoName(t *testing.T) {
	got := RepoName("witcomp1000-fall26", "a1", "alice")
	want := "witcomp1000-fall26-a1-alice"
	if got != want {
		t.Errorf("RepoName() = %q, want %q", got, want)
	}
}

func TestSaveAndLoadClassroom(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "CLASSROOM.yaml")

	c := &Classroom{
		Name:   "witcomp1000-fall26",
		Org:    "witcomp1000",
		Roster: "roster.csv",
	}

	if err := SaveClassroom(path, c); err != nil {
		t.Fatalf("SaveClassroom() error = %v", err)
	}

	got, err := LoadClassroom(path)
	if err != nil {
		t.Fatalf("LoadClassroom() error = %v", err)
	}

	if !reflect.DeepEqual(c, got) {
		t.Errorf("LoadClassroom() = %+v, want %+v", got, c)
	}
}

func TestLoadClassroom_MissingFile(t *testing.T) {
	_, err := LoadClassroom(filepath.Join(t.TempDir(), "nope.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestDefaultClassroomFile(t *testing.T) {
	if DefaultClassroomFile != "classroom.yaml" {
		t.Errorf("DefaultClassroomFile = %q, want %q", DefaultClassroomFile, "classroom.yaml")
	}
}
