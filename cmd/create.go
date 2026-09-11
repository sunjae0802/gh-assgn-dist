package cmd

import (
	"fmt"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
	"github.com/sunjae0802/gh-assgn-dist/internal"
)

var createName string
var createOrg string
var createRoster string
var createClassroom string
var createDryRun bool

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new classroom",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		roster := createRoster
		if roster == "" {
			roster = internal.DefaultRosterFile
		}

		// Verify the roster is readable before doing any network work, so a
		// typo'd path fails here rather than at the first distribute
		students, rosterWarnings, err := internal.LoadRoster(roster)
		if err != nil {
			return fmt.Errorf("could not read roster %q: %w", roster, err)
		}
		for _, w := range rosterWarnings {
			fmt.Printf("warning: %s\n", w)
		}
		if len(students) == 0 {
			return fmt.Errorf("roster %q has no students", roster)
		}

		client, err := api.DefaultRESTClient()
		if err != nil {
			return err
		}

		// Verify org exists
		var org struct{ Login string }
		if err := client.Get(fmt.Sprintf("orgs/%s", createOrg), &org); err != nil {
			return fmt.Errorf("org %q not found: %w", createOrg, err)
		}

		// Get authenticated user
		var user struct{ Login string }
		if err := client.Get("user", &user); err != nil {
			return err
		}

		// Verify user has admin access to org
		var membership struct{ Role string }
		if err := client.Get(fmt.Sprintf("orgs/%s/memberships/%s", createOrg, user.Login), &membership); err != nil {
			return fmt.Errorf("could not verify admin access to org %q: %w", createOrg, err)
		}
		if membership.Role != "admin" {
			return fmt.Errorf("user %q does not have admin access to org %q (role: %s)", user.Login, createOrg, membership.Role)
		}

		outFile := createClassroom
		if outFile == "" {
			outFile = internal.DefaultClassroomFile
		}
		if _, err := os.Stat(outFile); err == nil {
			return fmt.Errorf("%s already exists", outFile)
		}

		c := &internal.Classroom{
			Name:   createName,
			Org:    createOrg,
			Roster: roster,
		}

		if createDryRun {
			fmt.Printf("# would write %s (classroom: %s, org: %s, roster: %s with %s)\n",
				outFile, c.Name, c.Org, c.Roster, pluralStudents(len(students)))
			return nil
		}

		if err := internal.SaveClassroom(outFile, c); err != nil {
			return err
		}

		fmt.Printf("Created %s (%s in %s)\n", outFile, pluralStudents(len(students)), c.Roster)
		return nil
	},
}

// pluralStudents renders a student count with the right noun, e.g. "1 student".
func pluralStudents(n int) string {
	if n == 1 {
		return "1 student"
	}
	return fmt.Sprintf("%d students", n)
}

func init() {
	CreateCmd.Flags().StringVar(&createName, "name", "", "Classroom name, used as the student repo prefix (required)")
	CreateCmd.Flags().StringVar(&createOrg, "org", "", "GitHub organization name (required)")
	CreateCmd.Flags().StringVar(&createRoster, "roster", "", "Path to roster CSV file (default \"roster.csv\")")
	CreateCmd.Flags().StringVar(&createClassroom, "classroom", "", "Classroom YAML file to write (default \"classroom.yaml\")")
	CreateCmd.Flags().BoolVar(&createDryRun, "dry-run", false, "Verify the org, then print what would be written without writing it")
	CreateCmd.MarkFlagRequired("name")
	CreateCmd.MarkFlagRequired("org")
}
