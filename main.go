package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sunjae0802/gh-assgn-dist/cmd"
)

func main() {
	root := &cobra.Command{
		Use:   "assgn-dist",
		Short: "GitHub Classroom assignment distribution tool",
	}

	root.AddCommand(cmd.NewCmd)
	root.AddCommand(cmd.CreateCmd)
	root.AddCommand(cmd.CloneCmd)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
