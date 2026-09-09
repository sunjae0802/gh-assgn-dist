package internal

import (
	"encoding/csv"
	"fmt"
	"os"
)

// DefaultRosterFile is the roster file name used when --roster is not given.
const DefaultRosterFile = "roster.csv"

type Student struct {
	Name   string
	Email  string
	GitHub string
}

func LoadRoster(path string) ([]Student, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("roster file is empty")
	}

	// Skip header row if present
	start := 0
	if records[0][0] == "name" || records[0][0] == "Name" {
		start = 1
	}

	var students []Student
	for _, row := range records[start:] {
		if len(row) < 3 {
			continue
		}
		students = append(students, Student{
			Name:   row[0],
			Email:  row[1],
			GitHub: row[2],
		})
	}
	return students, nil
}
