package main

import (
	"fmt"
	"os"

	"github.com/locke23/git-clone/commands"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lit",
	Short: "lit is a minimalistic git cli clone",
	Long:  "lit is a minimalistic git cli clone",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func main() {
	rootCmd.AddCommand(commands.Init)
	rootCmd.AddCommand(commands.Add)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
