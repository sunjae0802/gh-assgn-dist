package cmd

import (
	"fmt"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
	"github.com/sunjae0802/gh-assgn-dist/internal"
)

var newOrg string
var newRoster string

var NewCmd = &cobra.Command{
	Use:   "new CLASSROOM",
	Short: "Create a new classroom",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		classroom := args[0]

		client, err := api.DefaultRESTClient()
		if err != nil {
			return err
		}

		// Verify org exists
		var org struct{ Login string }
		if err := client.Get(fmt.Sprintf("orgs/%s", newOrg), &org); err != nil {
			return fmt.Errorf("org %q not found: %w", newOrg, err)
		}

		// Get authenticated user
		var user struct{ Login string }
		if err := client.Get("user", &user); err != nil {
			return err
		}

		// Verify user has admin access to org
		var membership struct{ Role string }
		if err := client.Get(fmt.Sprintf("orgs/%s/memberships/%s", newOrg, user.Login), &membership); err != nil {
			return fmt.Errorf("could not verify admin access to org %q: %w", newOrg, err)
		}
		if membership.Role != "admin" {
			return fmt.Errorf("user %q does not have admin access to org %q (role: %s)", user.Login, newOrg, membership.Role)
		}

		outFile := classroom + ".yaml"
		if _, err := os.Stat(outFile); err == nil {
			return fmt.Errorf("%s already exists", outFile)
		}

		c := &internal.Classroom{
			Name:   classroom,
			Org:    newOrg,
			Roster: newRoster,
		}
		if err := internal.SaveClassroom(outFile, c); err != nil {
			return err
		}

		fmt.Printf("Created %s\n", outFile)
		return nil
	},
}

func init() {
	NewCmd.Flags().StringVar(&newOrg, "org", "", "GitHub organization name (required)")
	NewCmd.Flags().StringVar(&newRoster, "roster", "", "Path to roster CSV file (required)")
	NewCmd.MarkFlagRequired("org")
	NewCmd.MarkFlagRequired("roster")
}
