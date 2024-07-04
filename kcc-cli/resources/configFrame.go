package resources

import (
	"github.com/spf13/viper"
	"log"
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

// isKubeConfig check if a file is kube config
func isKubeConfig(filepath string) bool {
	v := viper.New()
	v.SetConfigFile(filepath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("ERROR: reading file from working directory %s", filepath)
		return false
	}
	var config KubeConfig
	if err := v.Unmarshal(&config); err != nil {
		log.Fatalf("ERROR: failed to unmarshal %s", filepath)
	}

	//check if required fields are present to qualify as kubeconfig
	if config.APIVersion == "" || config.Kind == "" || len(config.Clusters) == 0 || len(config.Contexts) == 0 {
		return false
	}
	return true
}
