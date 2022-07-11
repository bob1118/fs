/*
Copyright © 2022 bob

*/
package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bob1118/fs/db"
	"github.com/bob1118/fs/routers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// gatewayCmd represents the gateway command
var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "http server for freeswitch mod_xml_curl ",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

// run http gateway for fs mod_xml_curl.
fs gateway --run
// print gateway config
fs gateway
`,
	Run: gatewayCmdRun,
}

func init() {
	rootCmd.AddCommand(gatewayCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gatewayCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gatewayCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	gatewayCmd.Flags().BoolP("run", "r", false, "gateway run")
}

func gatewayCmdRun(cmd *cobra.Command, args []string) {
	fmt.Println("gateway called")
	if isRun, _ := cmd.Flags().GetBool(`run`); isRun {
		gatewayHttp()
	} else { //print config gateway
		gateway := viper.Sub(`gateway`)
		var list string
		allkeys := gateway.AllKeys()
		for _, key := range allkeys {
			list = fmt.Sprintf("%s\n%-30s=>%s", list, key, gateway.GetString(key))
		}
		log.Println(`gateway config:`, list)
	}
}

func gatewayHttp() {
	db.Initdb()
	h := routers.NewRouter()
	s := &http.Server{
		Addr:           viper.GetString(`gateway.http.addr`),
		Handler:        h,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
