/*
Copyright © 2024 asudarsanan@gmail.com
*/
package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kcc-cli/m/v2/resources"
	"kcc-cli/m/v2/ui"
	"log"
	"os"
	"path/filepath"
)

// rootCmd represents the base command when called without any subcommands
var (
	export  bool
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

// Selector This method passes the kubeconfig_path or kubeconfigs based on the user selected export flag, read from the
// config file. This methond invokes the UI selector modules and renders the selector.
func Selector(cmd *cobra.Command, args []string) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("ERROR: failed in getting user home directory %s", err)
	}

	configDir := filepath.Join(home, ".config", "kcc")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	configPath := filepath.Join(configDir, "config")
	viper.SetConfigFile(configPath)
	if !resources.FileExists(configPath) {
		log.Fatalf("ERROR: config file does not exist: %s", configPath)
	}

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("ERROR: failed to read the config file %s", err)
	}

	var kubeConfigPath string

	if !export {
		kubeConfigPath = viper.GetString("kubeconfig_path")
		err = processKubeConfig(kubeConfigPath)
		if err != nil {
			log.Printf("ERROR: failed to process kubeconfig: %v\n", err)
		}
	} else {
		kubeConfigPaths := viper.GetStringSlice("kubeconfigs")
		for _, path := range kubeConfigPaths {
			err = processKubeConfig(path)
			if err != nil {
				log.Printf("ERROR: failed to process kubeconfig: %v\n", err)
				return
			}
		}
	}
}

func processKubeConfig(kubeConfigPath string) error {
	config, err := resources.ReadKubeConfig(kubeConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	contexts := make([]string, len(config.Contexts))
	for i, ctx := range config.Contexts {
		contexts[i] = ctx.Name
	}
	contextPosition, contextList := ui.CussorPositionPointer(config)
	selectedContext, err := ui.ShowSelector(contextList, contextPosition)
	if err != nil {
		return fmt.Errorf("failed in selecting context: %w", err)
	}
	selected, err := ui.SwitchContext(config, selectedContext)
	if err != nil {
		return fmt.Errorf("failed in switching context: %w", err)
	}

	err = resources.WriteKubeConfig(kubeConfigPath, config)
	if err != nil {
		return fmt.Errorf("failed in writing kubeconfig: %w", err)
	}

	cyan := color.New(color.FgHiCyan).SprintFunc()
	fmt.Printf("Switched to context: %s\n", cyan(selected))
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
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
	rootCmd.PersistentFlags().BoolVarP(&export, "export", "e", false, "Export the selection as and environment variable. e.g; export KUBECONFIG=<selected context>")
	//rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(listCmd)
}
