package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Classroom struct {
	Name   string `yaml:"classroom"`
	Org    string `yaml:"org"`
	Roster string `yaml:"roster"`
}

func LoadClassroom(path string) (*Classroom, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Classroom
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func SaveClassroom(path string, c *Classroom) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// RepoName returns the student repo name for a given assignment, following the
// CLASSROOM-ASSGN-USERNAME convention.
func RepoName(classroom, assgn, github string) string {
	return fmt.Sprintf("%s-%s-%s", classroom, assgn, github)
}

// FindClassroomFile returns the single .yaml file in the current directory,
// or an error if there are zero or more than one.
func FindClassroomFile() (string, error) {
	matches, err := filepath.Glob("*.yaml")
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("no .yaml file found in current directory")
	}
	if len(matches) > 1 {
		return "", fmt.Errorf("multiple .yaml files found; specify one with --classroom")
	}
	return matches[0], nil
}
