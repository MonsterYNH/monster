/*
Copyright © 2021 yangniuhong

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package command

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var (
	loggerURI        string
	settingURI       string
	grpcEndpoint     string
	httpEndpoint     string
	httpEnable       *bool
	autoReloadConfig *bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "monster",
	Short: "go microservice util",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	rootCmd.Flags().StringVarP(&loggerURI, "logger", "l", "file://./logger.log", "")
	rootCmd.Flags().StringVarP(&settingURI, "setting", "s", "file://./config.ini", "")
	rootCmd.Flags().StringVarP(&grpcEndpoint, "grpc_endpoint", "g", "0.0.0.0:1234", "")
	rootCmd.Flags().StringVarP(&httpEndpoint, "http_endpoint", "t", "0.0.0.0:1235", "")
	httpEnable = rootCmd.Flags().BoolP("http_enable", "e", false, "")

	_, err := url.Parse(loggerURI)
	if err != nil {
		panic(fmt.Sprintf("ERROR: parse logger uri failed, error: %s", err.Error()))
	}
	setting, err := url.Parse(settingURI)
	if err != nil {
		panic(fmt.Sprintf("ERROR: parse setting uri failed, error: %s", err.Error()))
	}

	switch strings.ToLower(setting.Scheme) {
	case "file":
		viper.SetConfigFile(setting.Path)
	case "etcd":
		viper.AddRemoteProvider("etcd", fmt.Sprintf("%s:%s", setting.Host, setting.Port()), "/config/config.ini")
		viper.SetConfigType("ini")
	default:
		panic(fmt.Sprintf("ERROR: unknown option: %s", setting.Scheme))
	}

	// if cfgFile != "" {
	// 	// Use config file from the flag.
	// 	viper.SetConfigFile(cfgFile)
	// } else {
	// 	// Find home directory.
	// 	home, err := homedir.Dir()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		os.Exit(1)
	// 	}

	// 	// Search config in home directory with name ".monster" (without extension).
	// 	viper.AddConfigPath(home)
	// 	viper.SetConfigName(".monster")
	// }

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
