package resources

import (
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

// FileExists Check if the provided file exists in local machine
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// ReadKubeConfig Read the config file and form data structures.
func ReadKubeConfig(filePath string) (*KubeConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var config KubeConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// WriteKubeConfig Write back the config file - .kube/config
func WriteKubeConfig(filePath string, config *KubeConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

// ScanKubeConfigFiles scans for configuration files inside a provided directory path
// and validates if the identified file is a valid KubeConfig file.
func ScanKubeConfigFiles(dirPath string) ([]string, error) {
	var kubeConfigFiles []string

	entries, err := os.ReadDir(dirPath)
	fmt.Print(entries)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			path := filepath.Join(dirPath, entry.Name())
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if ext == ".yaml" || ext == ".yml" {
				if isKubeConfig(path) {
					kubeConfigFiles = append(kubeConfigFiles, path)
				}
			}
		}
	}

	return kubeConfigFiles, nil
}

func TrackThisConfigs(filePaths []string) {
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
	if FileExists(configPath) {
		err = viper.ReadInConfig()
		if err != nil {
			log.Fatalf("ERROR: failed to read the config file %s", err)
		}
		fmt.Print(viper.Get("kubeconfig_path"))
	}

}
