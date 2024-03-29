/*
Copyright © 2022 bob
*/
package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/bob1118/fs/fsconf"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// fsconfigCmd represents the fsconfig command
var fsconfigCmd = &cobra.Command{
	Use:   "fsconfig",
	Short: "freeswitch configuration, default dir is /etc/freeswitch",
	Long: `freeswitch configuration directory init and reset.
For example:
fs config fsconfig --init
fs config fsconfig --reset`,
	Run: fsconfigCmdRun,
}

func init() {
	configCmd.AddCommand(fsconfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fsconfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fsconfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	fsconfigCmd.Flags().BoolP("init", "i", false, "init bootable configurations")
	fsconfigCmd.Flags().BoolP("reset", "r", false, "reset bootable configurations to default")
	//fsconfigCmd.MarkFlagsMutuallyExclusive("init", "reset")
}

func fsconfigCmdRun(cmd *cobra.Command, args []string) {
	var dir string
	fmt.Println("fsconfig called")
	v := viper.GetViper()
	conf := v.GetString(`switch.conf`)
	if _, err := os.Stat(conf); err != nil { //conf from .fs switch.conf not exist.
		fmt.Println("switch.conf:", conf, "can not accessed")
		runos := runtime.GOOS
		switch runos {
		case `linux`:
			dir = `/etc/freeswitch`
		case `windows`:
			dir = `C:/Program Files/FreeSWITCH/conf`
		case `darwin`: //homebrew apple silinc
			dir = `/opt/homebrew/opt/freeswitch/etc/freeswitch`
		default:
		}
		if _, e := os.Stat(dir); e == nil { //default conf dir exist.
			v.Set(`switch.conf`, dir)
			v.WriteConfig()
			fmt.Println("set switch.conf to default:", dir)
		}
	} else {
		dir = conf
	}

	//--reset
	if isReset, _ := cmd.Flags().GetBool(`reset`); isReset {
		fmt.Println(fsconfigCmdReset(dir))
	}
	//--init
	if isInit, _ := cmd.Flags().GetBool(`init`); isInit {
		fmt.Println(fsconfigCmdInit(dir))
	}
}

func fsconfigCmdInit(dir string) error { return fsconf.NewConf(dir).Init() }

func fsconfigCmdReset(dir string) error { return fsconf.NewConf(dir).Reset() }
