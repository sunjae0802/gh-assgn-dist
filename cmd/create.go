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
			Roster: createRoster,
		}

		if createDryRun {
			fmt.Printf("# would write %s (classroom: %s, org: %s, roster: %s)\n",
				outFile, c.Name, c.Org, c.Roster)
			return nil
		}

		if err := internal.SaveClassroom(outFile, c); err != nil {
			return err
		}

		fmt.Printf("Created %s\n", outFile)
		return nil
	},
}

func init() {
	CreateCmd.Flags().StringVar(&createName, "name", "", "Classroom name, used as the student repo prefix (required)")
	CreateCmd.Flags().StringVar(&createOrg, "org", "", "GitHub organization name (required)")
	CreateCmd.Flags().StringVar(&createRoster, "roster", "", "Path to roster CSV file (required)")
	CreateCmd.Flags().StringVar(&createClassroom, "classroom", "", "Classroom YAML file to write (default \"classroom.yaml\")")
	CreateCmd.Flags().BoolVar(&createDryRun, "dry-run", false, "Verify the org, then print what would be written without writing it")
	CreateCmd.MarkFlagRequired("name")
	CreateCmd.MarkFlagRequired("org")
	CreateCmd.MarkFlagRequired("roster")
}
