/*
Copyright © 2022 bob

*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "fs cli toolset configuration, default filename is .fs",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

fs config --set switch.conf=/etc/freeswitch
fs config --set gateway.url=http://localhost/fsapi
fs config --get switch.conf`,
	Run: configCmdRun,
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	configCmd.Flags().BoolP("list", "l", false, "list all config options.")
	configCmd.Flags().String("set", "", "set config value, --set key=value")
	configCmd.Flags().String("get", "", "get config value, --get key return value")
	configCmd.MarkFlagsMutuallyExclusive("list", "set", "get")
}

func configCmdRun(cmd *cobra.Command, args []string) {
	fmt.Println("config called")
	//cmd.Flags().VisitAll(func(f *pflag.Flag) { fmt.Println(f.Name, f.Value) })

	//flag --list
	if isList, err := cmd.Flags().GetBool(`list`); err == nil {
		if isList {
			configCmdList()
		}
	}

	//flag --set
	//--set switch.db.host=127.0.0.1
	if set, err := cmd.Flags().GetString(`set`); err == nil {
		if len(set) > 0 {
			configCmdSetVar(set)
		}
	}

	//flag --get;
	//--get switch.db.host
	if get, err := cmd.Flags().GetString(`get`); err == nil {
		if len(get) > 0 {
			fmt.Println(get, "=>", configCmdGetVar(get))
		}
	}
}

func configCmdList() {
	fmt.Println(viper.GetViper().AllSettings())
}

func configCmdGetVar(key string) string { return viper.GetString(key) }

func configCmdSetVar(set string) {

	var isfound bool
	var key, value string

	spe := `=`
	kv := strings.ToLower(set)
	if strings.Contains(kv, spe) {
		if k, v, is := strings.Cut(kv, spe); is {
			key = strings.TrimSpace(k)
			value = strings.TrimSpace(v)
			isfound = is
		}
	}
	if isfound {
		if viper.IsSet(key) {
			viper.Set(key, value)
			viper.WriteConfig()
		} else {
			fmt.Println(`--set fail, undefine key:`, key)
		}
	}
}
