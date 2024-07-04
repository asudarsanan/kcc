package cmd

import (
	"github.com/spf13/cobra"
	"kcc-cli/m/v2/resources"
	"log"
	"strings"
)

var (
	wdPaths       string
	addDirPathCmd = &cobra.Command{
		Use:              "add-path",
		Short:            "Add the specified path for working directory look-up",
		Run:              addDirPathCommand,
		TraverseChildren: true,
	}
)

func addDirPathCommand(cmd *cobra.Command, args []string) {

	filePaths := strings.Split(wdPaths, ":")
	for _, filePath := range filePaths {
		configs, err := resources.ScanKubeConfigFiles(filePath)
		if err != nil {
			log.Fatalf("ERROR: failed to initialize working directory %s", err)
		}
		resources.TrackThisConfigs(configs)
	}
}

func init() {
	addDirPathCmd.Flags().StringVarP(&wdPaths, "path", "p", "", "Working directory paths for kube config lookup")
}
