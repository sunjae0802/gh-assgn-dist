package internal

import (
	"os"
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

func TestFindClassroomFile(t *testing.T) {
	t.Run("no yaml files", func(t *testing.T) {
		chdir(t, t.TempDir())
		if _, err := FindClassroomFile(); err == nil {
			t.Fatal("expected error when no .yaml files present, got nil")
		}
	})

	t.Run("single yaml file", func(t *testing.T) {
		dir := t.TempDir()
		chdir(t, dir)
		if err := os.WriteFile(filepath.Join(dir, "CLASSROOM.yaml"), []byte(""), 0644); err != nil {
			t.Fatal(err)
		}
		got, err := FindClassroomFile()
		if err != nil {
			t.Fatalf("FindClassroomFile() error = %v", err)
		}
		if got != "CLASSROOM.yaml" {
			t.Errorf("FindClassroomFile() = %q, want %q", got, "CLASSROOM.yaml")
		}
	})

	t.Run("multiple yaml files", func(t *testing.T) {
		dir := t.TempDir()
		chdir(t, dir)
		for _, name := range []string{"a.yaml", "b.yaml"} {
			if err := os.WriteFile(filepath.Join(dir, name), []byte(""), 0644); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := FindClassroomFile(); err == nil {
			t.Fatal("expected error for multiple .yaml files, got nil")
		}
	})
}

// chdir switches the working directory for the duration of the test and
// restores it afterward.
func chdir(t *testing.T, dir string) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(orig); err != nil {
			t.Fatal(err)
		}
	})
}
