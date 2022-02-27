package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "gominion",
	Short: "GoMinion is a CLI task manager",
	Long:  "A project I built for funsies.",
	// Run: func(cmd *cobra.Command, args []string){
	// 	// TODO
	// }
}
