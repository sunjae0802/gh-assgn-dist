package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/sunjae0802/gh-assgn-dist/internal"
)

var studentReposClassroom string

var StudentReposCmd = &cobra.Command{
	Use:   "student-repos ASSGN",
	Short: "Clone or update student repos for an assignment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		assgn := args[0]

		classroomFile := studentReposClassroom
		if classroomFile == "" {
			classroomFile = internal.DefaultClassroomFile
		}

		c, err := internal.LoadClassroom(classroomFile)
		if err != nil {
			return err
		}

		students, err := internal.LoadRoster(c.Roster)
		if err != nil {
			return err
		}

		for _, student := range students {
			repoName := internal.RepoName(c.Name, assgn, student.GitHub)
			localDir := repoName

			if _, err := os.Stat(localDir); os.IsNotExist(err) {
				fmt.Printf("Cloning %s/%s...\n", c.Org, repoName)
				out, err := exec.Command("gh", "repo", "clone", fmt.Sprintf("%s/%s", c.Org, repoName), localDir).CombinedOutput()
				if err != nil {
					fmt.Printf("warning: failed to clone %s: %v\n%s\n", repoName, err, out)
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

func init() {
	StudentReposCmd.Flags().StringVar(&studentReposClassroom, "classroom", "", "Classroom YAML file (default \"classroom.yaml\")")
}
