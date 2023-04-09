/*
Copyright © 2022 bob
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

const defaultFsContent = `#default is $HOME/.fs, format yaml.
gateway:
    db:
        host: 127.0.0.1
        name: freeswitch
        password: fsdba
        tableprefix: g
        user: fsdba
    enablea1hash: false
    http:
        addr: localhost:8080
postgres:
    host: 127.0.0.1
    name: postgres
    password: postgres
    user: postgres
server:
    db:
        host: 127.0.0.1
        name: freeswitch
        password: fsdba
        tableprefix: s
        user: fsdba
    http:
        addr: 10.10.10.21:8080
        readtimeout: 4
        writetimeout: 4
    eventsocket:
        serveraddr: 127.0.0.1:12345
switch:
    conf: /etc/freeswitch
    cdr:
        modname: mod_odbc_cdr
        a-leg: cdr_table_a_leg
        b-leg: cdr_table_b_leg
        both: 
    db:
        host: 127.0.0.1
        name: freeswitch
        password: fsdba
        user: fsdba
    eventsocket:
        ipaddr: 127.0.0.1
        port: 8021
        password: ClueCon
    record:
        dir: /var/lib/freeswitch/recorddings
    vars:
        ipv4: 10.10.10.21
        external_sip_ip: $${local_ip_v4}
        external_rtp_ip: $${local_ip_v4}
    xml_curl:
        url: http://localhost:8080/fsapi
        bindings: dialplan|configuration|directory|phrases
`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fs",
	Short: "switch command line toolset",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. 

For example:
fs config --set switch.conf=/etc/freeswitch
fs config --get switch.xml_curl.url
http://localhost/fsapi

fs config fsconfig --reset --init

fs gateway --run
fs server --run
systemctl restart freeswitch
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: rootCmdRun,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fs)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		var configFile string
		filename := `.fs`
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".fs" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(filename)
		viper.SetConfigType("yaml")

		//if defaultFile not exist, write defaultFsContent to it.
		configFile = filepath.Join(home, filename)
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			if err := os.WriteFile(configFile, []byte(defaultFsContent), 0644); err != nil {
				fmt.Println(err)
			}
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("viper.ReadInConfig()", err)
	}
}

func rootCmdRun(cmd *cobra.Command, args []string) {
	fmt.Println(`rootCmd called`)
}
