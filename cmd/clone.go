package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sunjae0802/gh-assgn-dist/internal"
)

var cloneClassroom string
var cloneDryRun bool

var CloneCmd = &cobra.Command{
	Use:   "clone ASSGN",
	Short: "Clone or update student repos for an assignment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		assgn := args[0]

		classroomFile := cloneClassroom
		if classroomFile == "" {
			classroomFile = internal.DefaultClassroomFile
		}

		c, err := internal.LoadClassroom(classroomFile)
		if err != nil {
			return err
		}

		students, rosterWarnings, err := internal.LoadRoster(c.Roster)
		if err != nil {
			return err
		}
		for _, w := range rosterWarnings {
			fmt.Printf("warning: %s\n", w)
		}

		for _, student := range students {
			repoName := internal.RepoName(c.Name, assgn, student.GitHub)
			localDir := filepath.Join(assgn, repoName)

			_, statErr := os.Stat(localDir)
			missing := os.IsNotExist(statErr)

			if cloneDryRun {
				if missing {
					fmt.Printf("gh repo clone %s/%s %s\n", c.Org, repoName, localDir)
				} else {
					fmt.Printf("git -C %s pull\n", localDir)
				}
				continue
			}

			if missing {
				fmt.Printf("Cloning %s/%s...\n", c.Org, repoName)
				out, err := exec.Command("gh", "repo", "clone", fmt.Sprintf("%s/%s", c.Org, repoName), localDir).CombinedOutput()
				if err != nil {
					fmt.Println(cloneWarning(c.Org, repoName, out))
				}
			} else {
				fmt.Printf("Pulling %s...\n", repoName)
				out, err := exec.Command("git", "-C", localDir, "pull").CombinedOutput()
				if err != nil {
					fmt.Printf("warning: failed to pull %s: %v\n%s\n", repoName, err, out)
				}
			}
		}

		return nil
	},
}

// cloneWarning turns gh's raw clone output into a short, friendly warning.
func cloneWarning(org, repoName string, out []byte) string {
	if strings.Contains(string(out), "Could not resolve to a Repository") {
		return fmt.Sprintf("warning: repository `%s/%s` not found", org, repoName)
	}
	return fmt.Sprintf("warning: failed to clone %s/%s: %s", org, repoName, strings.TrimSpace(string(out)))
}

func init() {
	CloneCmd.Flags().StringVar(&cloneClassroom, "classroom", "", "Classroom YAML file (default \"classroom.yaml\")")
	CloneCmd.Flags().BoolVar(&cloneDryRun, "dry-run", false, "Print the git/gh commands without executing")
}
