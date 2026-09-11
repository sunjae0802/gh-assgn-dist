package internal

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

// DefaultRosterFile is the roster file name used when --roster is not given.
const DefaultRosterFile = "roster.csv"

type Student struct {
	Name   string
	Email  string
	GitHub string
}

// LoadRoster reads a roster CSV. The file must start with a header row naming a
// GitHub column (`github` in any capitalization); `name` and `email` columns are
// used when present. Columns may appear in any order.
//
// Blank rows are skipped silently. A row carrying a name or an email but no
// GitHub username is skipped too, and reported in the returned warnings — it is
// usually a roster typo rather than a reason to abort the whole run.
func LoadRoster(path string) ([]Student, []string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	// Rows shorter or longer than the header are tolerated here so that a
	// single malformed line can be reported as a warning instead of failing
	// the whole roster.
	r.FieldsPerRecord = -1

	header, err := r.Read()
	if err == io.EOF {
		return nil, nil, fmt.Errorf("roster file is empty")
	}
	if err != nil {
		return nil, nil, err
	}

	// Spreadsheets often save UTF-8 with a byte-order mark, which would
	// otherwise hide in the first header name.
	if len(header) > 0 {
		header[0] = strings.TrimPrefix(header[0], "\ufeff")
	}

	col := map[string]int{}
	for i, name := range header {
		key := strings.ToLower(strings.TrimSpace(name))
		if _, seen := col[key]; !seen {
			col[key] = i
		}
	}

	githubCol, ok := col["github"]
	if !ok {
		return nil, nil, fmt.Errorf("roster is missing a %q column (header: %s)",
			"github", strings.Join(header, ","))
	}
	nameCol, hasName := col["name"]
	emailCol, hasEmail := col["email"]

	field := func(row []string, i int, present bool) string {
		if !present || i >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[i])
	}

	var students []Student
	var warnings []string
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}

		// FieldPos reports the line the row actually came from, which is not
		// the row's ordinal once blank lines are skipped.
		line, _ := r.FieldPos(0)

		student := Student{
			Name:   field(row, nameCol, hasName),
			Email:  field(row, emailCol, hasEmail),
			GitHub: field(row, githubCol, true),
		}

		if student.GitHub != "" {
			students = append(students, student)
			continue
		}

		// No username: silent for a blank line, a warning when the row
		// names someone.
		if student.Name == "" && student.Email == "" {
			continue
		}
		warnings = append(warnings, fmt.Sprintf(
			"line %d of %s has no GitHub username, skipping (%s)",
			line, path, describeStudent(student)))
	}

	return students, warnings, nil
}

// describeStudent labels a row in a warning using whichever of name and email
// the row actually filled in.
func describeStudent(s Student) string {
	switch {
	case s.Name != "" && s.Email != "":
		return fmt.Sprintf("%s <%s>", s.Name, s.Email)
	case s.Name != "":
		return s.Name
	default:
		return s.Email
	}
}
