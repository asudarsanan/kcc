/*
Copyright © 2024 asudarsanan@gmail.com
*/
package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"kcc-cli/m/v2/resources"
	"kcc-cli/m/v2/ui"
	"log"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "kcc",
		Short: "A Kubernetes Context Controller",
		//Long: `A longer description that spans multiple lines and likely contains
		//examples and usage of using your application. For example:
		//
		//Cobra is a CLI library for Go that empowers applications.
		//This application is a tool to generate the needed files
		//to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: Selector,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.

func Selector(cmd *cobra.Command, args []string) {
	kubeConfigPath := resources.KubeConfigPath()
	config, err := resources.ReadKubeConfig(kubeConfigPath)
	if err != nil {
		log.Printf("ERROR: failed to read kubeconfig: %v\n", err)
		return
	}
	contexts := make([]string, len(config.Contexts))
	for i, ctx := range config.Contexts {
		contexts[i] = ctx.Name
	}
	contextPosition, contextList := ui.CussorPositionPointer(config)
	selectedContext, err := ui.ShowSelector(contextList, contextPosition)
	if err != nil {
		log.Printf("ERROR: failed in selecting context: %v\n", err)
		return
	}
	selected, err := ui.SwitchContext(config, selectedContext)
	if err != nil {
		log.Printf("ERROR: failed in switching context: %v\n", err)
		return
	}
	err = resources.WriteKubeConfig(kubeConfigPath, config)
	if err != nil {
		log.Printf("ERROR: failed in writing kubeconfig: %v\n", err)
	}
	cyan := color.New(color.FgHiCyan).SprintFunc()
	fmt.Printf("Switched to context: %s\n", cyan(selected))
}
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.v2.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	//rootCmd.AddCommand(listCmd)
}
