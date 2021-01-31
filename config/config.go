package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	rootCmd *cobra.Command

	configs = map[string]*Config{
		"helloworld_http": {
			ConfigURI: "/config.yaml",
			LoggerURI: "/log",
			Endpoint:  "0.0.0.0:1235",
			SubServeInfos: map[string]string{
				"helloworld_grpc": "0.0.0.0:1234",
			},
		},
		"helloworld_grpc": {
			ConfigURI: "/config.yaml",
			LoggerURI: "/log",
			Endpoint:  "0.0.0.0:1234",
		},
	}
)

// Config serve starup config
type Config struct {
	ConfigURI     string
	LoggerURI     string
	Endpoint      string
	SubServeInfos map[string]string
	Disable       bool
	ServeType     string
	Ext           map[string]string
}

// GetConfigEntry get serve start up config
func GetConfigEntry(name string) (*Config, error) {
	config, exist := configs[name]
	if !exist {
		return nil, fmt.Errorf("can not serve config, named: %s", name)
	}
	return config, nil
}
