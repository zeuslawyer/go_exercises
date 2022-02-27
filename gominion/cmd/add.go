package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{ // note it's unexported.
	Use:   "add",
	Short: "Adds a task to your task list",
	Run: func(cmd *cobra.Command, args []string) { // args are parsed as space separated
		// merge the args into a single string
		task := strings.Join(args, " ")
		fmt.Println("Added task: ", task)

	},
}

func init() {
	RootCmd.AddCommand(addCmd) // Mount the add command.
}
