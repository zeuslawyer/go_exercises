package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Marks a task as complete",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("do called")
		var taskIds []int

		for _, arg := range args {
			argInt, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("failed to parse the argument:", argInt, "- it isn't an integer.")
				os.Exit(1)
			} else {
				fmt.Println("Arg as int is: ", argInt)
				taskIds = append(taskIds, argInt)
			}
		}
		fmt.Println("here are your task ids: ", taskIds)
	},
}

func init() {
	RootCmd.AddCommand(doCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// doCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// doCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
