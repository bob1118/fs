/*
Copyright © 2022 bob

*/
package cmd

import (
	"fmt"
	"log"
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
        tableprefix: cc
        user: fsdba
    enablea1hash: false
    http:
        addr: localhost
    xml_curl:
        url: http://localhost/fsapi
        bindings: dialplan|configuration|directory|phrases
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
        user: fsdba
    http:
        addr: 10.10.10.10:80
        readtimeout: 4
        writetimeout: 4
    inbound:
        addr: 127.0.0.1:8021
        password: ClueCon
    outbound:
        addr: 127.0.0.1:12345
switch:
    conf: /etc/freeswitch
    db:
        host: 127.0.0.1
        name: freeswitch
        password: fsdba
        user: fsdba
`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fs",
	Short: "switch command line toolset",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. 

For example:
where is app configuration file ?
global flag "--config=where/file.format" set fs configuration file. 
if no flag "--config", fs search $HOME/.fs as a default.
how to set?
fs config --set switch.conf=/etc/freeswitch
how to get?
fs config --get switch.conf
//////fs config fsconfg --reset --init//////
how to run switch gateway?
fs gateway
hot to run switch server?
fs server
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
				log.Println(err)
			}
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		log.Println("viper.ReadInConfig()", err)
	}
}

func rootCmdRun(cmd *cobra.Command, args []string) {
	//cmd.Flags().VisitAll(func(f *pflag.Flag) {log.Println(f.Name, f.Value)})
	fmt.Println(`rootCmd called`)
}
