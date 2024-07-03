/*
Copyright © 2024 asudarsanan@gmail.com
*/
package cmd

import (
	"github.com/spf13/cobra"
	"kcc-cli/m/v2/cmd/list"
	"kcc-cli/m/v2/resources"
)

// lsCmd represents the ls command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all imported kube config contexts",
	Run:   listCommand,
}

func listCommand(cmd *cobra.Command, args []string) {
	contexts := resources.FetchContexts()
	list.Contexts(contexts)
}
func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//listCmd.Flags().BoolP("list", "l", false, "Help message for toggle")
}
