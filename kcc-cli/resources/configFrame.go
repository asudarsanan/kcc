package resources

import (
	"github.com/spf13/viper"
	"log"
)

type KubeConfig struct {
	APIVersion     string      `yaml:"apiVersion"`
	Clusters       []Cluster   `yaml:"clusters"`
	Contexts       []Context   `yaml:"contexts"`
	CurrentContext string      `yaml:"current-context"`
	Kind           string      `yaml:"kind"`
	Preferences    Preferences `yaml:"preferences"`
	Users          []User      `yaml:"users"`
}

type Cluster struct {
	Name    string        `yaml:"name"`
	Cluster ClusterDetail `yaml:"cluster"`
}

type ClusterDetail struct {
	Server                   string `yaml:"server"`
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
}

type Context struct {
	Name    string        `yaml:"name"`
	Context ContextDetail `yaml:"context"`
}

type ContextDetail struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

type User struct {
	Name string     `yaml:"name"`
	User UserDetail `yaml:"user"`
}

type UserDetail struct {
	ClientCertificateData string `yaml:"client-certificate-data"`
	ClientKeyData         string `yaml:"client-key-data"`
}

type Preferences struct {
	Colors bool `yaml:"colors,omitempty"`
}

type ContextWithClusterInfo struct {
	Name    string
	Context ContextDetail
	Cluster ClusterDetail
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
