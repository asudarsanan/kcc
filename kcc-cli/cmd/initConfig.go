package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kcc-cli/m/v2/resources"
	"log"
	"os"
	"path/filepath"
)

var (
	resetFlag bool
	initCmd   = &cobra.Command{
		Use:   "init",
		Short: "Initialize the KCC configuration",
		Run:   initConfigCmd,
		//TraverseChildren: true,
	}
)

func initConfigCmd(cmd *cobra.Command, args []string) {

	if resetFlag {

		err := removeSwap()
		if err != nil {
			log.Fatal("ERROR: failed to remove the swap", err)
		}
		viper.Set("kubeconfigs", "")
		viper.WriteConfig()
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")

		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("ERROR: failed in getting user home directory %s", err)
			//os.Exit(1)
		}
		configDir := filepath.Join(home, ".config", "kcc")
		configPath := filepath.Join(configDir, "config")
		//viper.SetConfigFile(configPath)
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			log.Fatalf("ERROR: failed in creating config dir %s", err)
		}
		viper.Set("default_config_reference", filepath.Join(configDir, "default-config.kcc"))
		viper.Set("kubeconfig_path", resources.KubeConfigPath())
		err = resources.CopyFile(resources.KubeConfigPath(), filepath.Join(configDir, "default-config.kcc"))
		if err != nil {
			log.Fatalf("ERROR: failed in copying default config file %s", err)
		}
		viper.SetDefault("kubeconfigs", "")
		viper.Set("active_swaps", "")
		viper.Set("swapped_file", "")

		if !resources.FileExists(configPath) {
			err = viper.WriteConfigAs(configPath)
			if err != nil {
				log.Fatalf("ERROR: failed in writing config file %s", err)
			}
			fmt.Println("Configuration initialized", viper.ConfigFileUsed())
		} else {
			fmt.Println("Configuration already initialized", resources.KccConfigFile())
			fmt.Println("You need to re-initialize the KCC configuration, using kcc init --clean-slate")
		}
	}

}

func removeSwap() error {
	configPath := resources.KccConfigFile()
	if resources.FileExists(configPath) {
		kccConfig := resources.KccConfigFile()
		viper.SetConfigFile(kccConfig)
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("ERROR: failed to read config file: %s", err)
			return err
		}

		// Check if there is a swapped_file to revert
		swappedFile := viper.GetString("swapped_file")
		defaultConfig := viper.GetString("default_config_reference")
		if swappedFile != "" && resources.FileExists(defaultConfig) {
			err = resources.CopyFile(defaultConfig, resources.KubeConfigPath())
			if err != nil {
				log.Fatalf("ERROR: failed to revert the kubeconfig file: %v", err)
				return err
			}
			// Clear the swapped_file key
			viper.Set("swapped_file", "")
			viper.Set("active_swaps", "")
			viper.WriteConfig()
			fmt.Println("Configuration reset successfully.")
		} else {
			fmt.Println("No swapped configuration to revert.")
		}
	} else {
		fmt.Println("Configuration file does not exist.")
	}
	return nil
}

func init() {
	initCmd.Flags().BoolVar(&resetFlag, "clean-slate", false, "Create a fresh configuration, overwriting existing one")
	initCmd.AddCommand(addDirPathCmd)
}
