package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v2"
)

type KubeConfig struct {
	APIVersion     string                 `yaml:"apiVersion"`
	Clusters       []Cluster              `yaml:"clusters"`
	Contexts       []Context              `yaml:"contexts"`
	CurrentContext string                 `yaml:"current-context"`
	Kind           string                 `yaml:"kind"`
	Preferences    map[string]interface{} `yaml:"preferences"`
	Users          []User                 `yaml:"users"`
}

type Cluster struct {
	Name    string                 `yaml:"name"`
	Cluster map[string]interface{} `yaml:"cluster"`
}

type Context struct {
	Name    string                 `yaml:"name"`
	Context map[string]interface{} `yaml:"context"`
}

type User struct {
	Name string                 `yaml:"name"`
	User map[string]interface{} `yaml:"user"`
}

// Read the config file and form data structures.
func readKubeConfig(filePath string) (*KubeConfig, error) {
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

// Write back the config file - .kube/config
func writeKubeConfig(filePath string, config *KubeConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// Switch the current-context based on the selection made
func switchContext(config *KubeConfig, contextName string) error {
	for _, context := range config.Contexts {
		if context.Name == contextName {
			config.CurrentContext = contextName
			return nil
		}
	}
	return fmt.Errorf("context %s not found", contextName)
}

// render selector
func showSelector(options []string) (string, error) {
	prompt := promptui.Select{
		Label: "Select Kubernetes cluster context:",
		Items: options,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

func main() {

	var kubeConfigPath string
	switch runtime.GOOS {
	case "windows":
		kubeConfigPath = filepath.Join(os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"), ".kube", "config")
	case "darwin", "linux":
		kubeConfigPath = os.Getenv("HOME") + "/.kube/config"
	default:
		log.Fatalf("Unsupported OS: %s", runtime.GOOS)
	}

	config, err := readKubeConfig(kubeConfigPath)
	if err != nil {
		fmt.Printf("Error reading kubeconfig: %v\n", err)
		return
	}

	contexts := make([]string, len(config.Contexts))
	for i, ctx := range config.Contexts {
		contexts[i] = ctx.Name
	}

	selectedContext, err := showSelector(contexts)
	if err != nil {
		fmt.Printf("Error selecting context: %v\n", err)
		return
	}

	err = switchContext(config, selectedContext)
	if err != nil {
		fmt.Printf("Error switching context: %v\n", err)
		return
	}

	err = writeKubeConfig(kubeConfigPath, config)
	if err != nil {
		fmt.Printf("Error writing kubeconfig: %v\n", err)
	}

	fmt.Printf("Switched to context %s\n", selectedContext)
}
