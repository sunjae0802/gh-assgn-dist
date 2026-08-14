package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
	"github.com/sunjae0802/gh-assgn-dist/internal"
)

var createTemplate string
var createClassroom string
var createDryRun bool

var CreateCmd = &cobra.Command{
	Use:   "create ASSGN",
	Short: "Create student repos for an assignment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		assgn := args[0]

		classroomFile := createClassroom
		if classroomFile == "" {
			var err error
			classroomFile, err = internal.FindClassroomFile()
			if err != nil {
				return err
			}
		}

		c, err := internal.LoadClassroom(classroomFile)
		if err != nil {
			return err
		}

		students, err := internal.LoadRoster(c.Roster)
		if err != nil {
			return err
		}

		template := createTemplate
		if template == "" {
			template = assgn
		}

		client, err := api.DefaultRESTClient()
		if err != nil {
			return err
		}

		// Verify template repo exists
		templateOwner, templateName, err := splitRepo(template)
		if err != nil {
			return err
		}
		var templateRepo struct{ FullName string `json:"full_name"` }
		if err := client.Get(fmt.Sprintf("repos/%s/%s", templateOwner, templateName), &templateRepo); err != nil {
			return fmt.Errorf("template repo %q not found: %w", template, err)
		}

		for _, student := range students {
			repoName := internal.RepoName(c.Name, assgn, student.GitHub)
			org := c.Org

			exists, err := repoExists(client, org, repoName)
			if err != nil {
				fmt.Printf("warning: failed to check if repo %s exists: %v\n", repoName, err)
				continue
			}

			if createDryRun {
				if exists {
					fmt.Printf("# %s/%s already exists, skipping creation\n", org, repoName)
				} else {
					fmt.Printf("gh api -X POST /repos/%s/%s/generate -f owner=%s -f name=%s -f private=true\n",
						templateOwner, templateName, org, repoName)
				}
				fmt.Printf("gh api -X PUT /repos/%s/%s/collaborators/%s -f permission=write\n",
					org, repoName, student.GitHub)
				continue
			}

			if exists {
				fmt.Printf("%s already exists, skipping creation\n", repoName)
			} else {
				// Create repo from template
				body := map[string]interface{}{
					"owner":   org,
					"name":    repoName,
					"private": true,
				}
				bodyBytes, _ := json.Marshal(body)
				var createdRepo struct{ FullName string `json:"full_name"` }
				if err := client.Post(fmt.Sprintf("repos/%s/%s/generate", templateOwner, templateName), bytes.NewReader(bodyBytes), &createdRepo); err != nil {
					fmt.Printf("warning: failed to create repo %s: %v\n", repoName, err)
					continue
				}
				fmt.Printf("Created %s\n", createdRepo.FullName)
			}

			// Add student as outside collaborator with write access (idempotent,
			// so this also repairs repos where a prior run's invite failed)
			collabBody := map[string]string{"permission": "write"}
			collabBytes, _ := json.Marshal(collabBody)
			var collabResp struct{}
			if err := client.Put(fmt.Sprintf("repos/%s/%s/collaborators/%s", org, repoName, student.GitHub), bytes.NewReader(collabBytes), &collabResp); err != nil {
				fmt.Printf("warning: failed to add collaborator %s to %s: %v\n", student.GitHub, repoName, err)
			} else {
				fmt.Printf("Added %s as collaborator on %s\n", student.GitHub, repoName)
			}
		}

		return nil
	},
}

func splitRepo(repo string) (string, string, error) {
	owner, name, ok := strings.Cut(repo, "/")
	if !ok {
		return "", "", fmt.Errorf("invalid repo format %q: expected owner/name", repo)
	}
	return owner, name, nil
}

// repoExists reports whether org/name already exists on GitHub.
func repoExists(client *api.RESTClient, org, name string) (bool, error) {
	var repo struct{ FullName string `json:"full_name"` }
	err := client.Get(fmt.Sprintf("repos/%s/%s", org, name), &repo)
	if err == nil {
		return true, nil
	}
	var httpErr *api.HTTPError
	if errors.As(err, &httpErr) && httpErr.StatusCode == 404 {
		return false, nil
	}
	return false, err
}

func init() {
	CreateCmd.Flags().StringVar(&createTemplate, "template", "", "Template repo (owner/name); defaults to assignment name")
	CreateCmd.Flags().StringVar(&createClassroom, "classroom", "", "Classroom YAML file; defaults to single .yaml in CWD")
	CreateCmd.Flags().BoolVar(&createDryRun, "dry-run", false, "Print gh api commands without executing")
}
