package internal

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// DefaultClassroomFile is the classroom file name used when --classroom is
// not given.
const DefaultClassroomFile = "classroom.yaml"

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
