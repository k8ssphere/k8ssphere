package config

import (
	"fmt"
	"github.com/spf13/viper"
	"k8ssphere.io/k8ssphere/pkg/simple/client/k8s/options"
)

const (
	// config_name is the default name of configuration
	config_name = "config"

	// config_path the default location of the configuration file
	config_path = "/etc/k8s"
)

type Config struct {
	KubernetesOptions *options.KubernetesOptions
}

/**

 */
func NewConfig() *Config {
	s := Config{
		KubernetesOptions: options.NewKubernatesOptions(),
	}
	return &s
}

func TryLoadFromDisk() (*Config, error) {
	viper.SetConfigName(config_name)
	viper.AddConfigPath(config_path)
	// Load from current working directory, only used for debugging
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, err
		} else {
			return nil, fmt.Errorf("error parsing configuration file %s", err)
		}
	}
	conf := NewConfig()
	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}
	return conf, nil
}
