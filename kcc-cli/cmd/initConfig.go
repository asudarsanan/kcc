package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	}
)

func initConfigCmd(cmd *cobra.Command, args []string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("ERROR: failed in getting user home directory %s", err)
		os.Exit(1)
	}
	configDir := filepath.Join(home, ".config", "kcc")
	configPath := filepath.Join(configDir, "config")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		log.Fatalf("ERROR: failed in creating config dir %s", err)
	}
	viper.SetDefault("kubeconfig_path", filepath.Join(home, ".kube", "config"))

	/* TODO we need to make a way to include a flag input for taking in working dir path from the user and append them if multivalues
	this means - --working-dir /home/test,/home/test2
	there needs to be a way to idenfity valid kubeconfig files inside these working dir and make a note of them in the config itself?
	this needs to happen at the time of initialization.
	*/

	if resetFlag {
		err = viper.WriteConfigAs(configPath)
		if err != nil {
			log.Fatalf("ERROR: failed in writing config file %s", err)
			os.Exit(1)
		}
		fmt.Println("Configuration initialized", viper.ConfigFileUsed())
	} else {
		fmt.Println("Configuration already initialized", viper.ConfigFileUsed())
		fmt.Println("You need to re-initialize the KCC configuration, using kcc init --clean-slate")
	}

}

func init() {
	initCmd.Flags().BoolVar(&resetFlag, "clean-slate", false, "Create a fresh configuration, overwriting existing one")
	rootCmd.AddCommand(initCmd)
}
