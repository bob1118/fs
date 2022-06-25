/*
Copyright © 2022 bob

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// fsconfigCmd represents the fsconfig command
var fsconfigCmd = &cobra.Command{
	Use:   "fsconfig",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fsconfig called")
	},
}

func init() {
	configCmd.AddCommand(fsconfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fsconfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fsconfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
