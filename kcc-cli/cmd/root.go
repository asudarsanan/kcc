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
	extend        bool
	defaultConfig bool
	rootCmd       = &cobra.Command{
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
	var kubeConfigPath []string
	if defaultConfig {
		err := removeSwap()
		if err != nil {
			log.Fatal("ERROR: failed in reverting to default configuration", err)
		}
	}
	if !extend {
		kubeConfigPath = viper.GetStringSlice("kubeconfig_path")
		err = processKubeConfig(kubeConfigPath)
		if err != nil {
			log.Printf("ERROR: failed to process kubeconfig: %v\n", err)
		}
	} else {
		kubeConfigPaths := viper.GetStringSlice("kubeconfigs")
		err = processKubeConfig(kubeConfigPaths)
		//for _, path := range kubeConfigPaths {
		//	err = processKubeConfig(path)
		if err != nil {
			log.Printf("ERROR: failed to process kubeconfig: %v\n", err)
			return
		}
		//}
	}
}

// processKubeConfig processes the provided kubeconfig paths to allow the user to select
// and switch contexts. If the export flag is true, it first lets the user select a config
// file from the provided paths before proceeding to context selection.
//
// Parameters:
// - kubeConfigPaths: A slice of strings containing paths to kubeconfig files.
// - export: A boolean flag indicating whether to export the selection as an environment variable.
//
// Returns:
// - error: An error object if any issues occur during the process.
//
// This function performs the following steps:
//  1. If the export flag is true, it displays a selector for the user to choose a kubeconfig file
//     from the provided paths. The selected path is then used for further processing.
//  2. Reads the kubeconfig file from the selected or default path and parses its contexts.
//  3. Creates a list of context names and displays a selector for the user to choose a context.
//  4. Switches to the selected context and updates the CurrentContext field in the kubeconfig.
//  5. Writes the updated kubeconfig back to the file.
//  6. Prints a confirmation message indicating the selected context.
func processKubeConfig(kubeConfigPaths []string) error {
	//var combinedContexts []resources.Context
	//var configs []*resources.KubeConfig
	var contextNames []string
	if extend {
		selectedConfigPath, err := ui.ShowPathSelector(kubeConfigPaths)
		if err != nil {
			return fmt.Errorf("failed in selecting config path: %w", err)
		}
		kubeConfigPaths = []string{selectedConfigPath}
		//log.Println("INFO: Selected config path:", selectedConfigPath)
	}
	// Assuming only one config path after the selection process
	configPath := kubeConfigPaths[0]
	config, err := resources.ReadKubeConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to read kubeconfig: %w", err)
	}
	contexts := config.Contexts
	contextNames = make([]string, len(contexts))
	for i, ctx := range contexts {
		contextNames[i] = ctx.Name
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
	config.CurrentContext = selectedContext
	err = resources.WriteKubeConfig(configPath, config)
	if err != nil {
		return fmt.Errorf("failed in writing kubeconfig: %w", err)
	}
	cyan := color.New(color.FgHiCyan).SprintFunc()
	fmt.Printf("Switched to context: %s\n", cyan(selected))
	if extend {
		err := resources.SwapThisConfig(kubeConfigPaths[0])
		if err != nil {
			return fmt.Errorf("failed in swapping this context: %w", err)
		}
	}
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
	rootCmd.PersistentFlags().BoolVarP(&extend, "extend", "e", false, "Make an external configuration the primary config and use with kubectl commands.")
	rootCmd.PersistentFlags().BoolVarP(&defaultConfig, "remove-swap", "r", false, "This will remove any active swaps in the config files and make you default kube config in $HOME/.kube/config")
	rootCmd.MarkFlagsMutuallyExclusive("extend", "remove-swap")
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(listCmd)
}
