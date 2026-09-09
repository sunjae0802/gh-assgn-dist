package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/sunjae0802/gh-assgn-dist/cmd"
)

func main() {
	root := &cobra.Command{
		Use:   "assgn-dist",
		Short: "GitHub Classroom assignment distribution tool",
	}

	// Usage is worth printing when the command line itself is wrong, but not
	// when a command fails while running. SilenceUsage is read after RunE
	// returns, so setting it here leaves flag parsing alone and silences only
	// runtime failures.
	root.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		cmd.SilenceUsage = true
	}

	root.AddCommand(cmd.CreateCmd)
	root.AddCommand(cmd.DistributeCmd)
	root.AddCommand(cmd.CloneCmd)

	// cobra prints the error to stderr itself; just set the exit status.
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
