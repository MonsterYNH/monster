package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	initConfigCmd.AddCommand(loggerCmd, configCmd)
}

var initConfigCmd = &cobra.Command{
	Use:   "init_config",
	Short: "config logger URI and config URI",
	Long:  "config format [schema]://[username]:[password]@[address]:[port]/[path]?[query]",
}

var loggerCmd = &cobra.Command{
	Use:   "logger",
	Short: "logger URI",
	Long:  "example: file:///app/config.ini?auto_save=true | etcd://monster:password@monster.ynh:2379",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("logger uri: ", args)
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config URI",
	Long:  "example: file:///app/config.ini?auto_save=true | etcd://monster:password@monster.ynh:2379",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config uri: ", args)
	},
}
