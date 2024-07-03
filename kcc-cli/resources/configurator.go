package resources

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func KubeConfigPath() string {
	var kubeConfigPath string
	switch runtime.GOOS {
	case "windows":
		kubeConfigPath = filepath.Join(os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"), ".kube", "config")
	case "darwin", "linux":
		kubeConfigPath = os.Getenv("HOME") + "/.kube/config"
	default:
		log.Fatalf("Unsupported OS: %s", runtime.GOOS)
	}
	return kubeConfigPath
}
func FetchContexts() []Context {
	config, err := ReadKubeConfig(KubeConfigPath())
	if err != nil {
		fmt.Printf("Error in loading kubeconfig: %v\n", err)
	}
	contexts := config.Contexts
	return contexts
}
