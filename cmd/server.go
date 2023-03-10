/*
Copyright © 2022 bob
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bob1118/fs/db"
	"github.com/bob1118/fs/esl"
	"github.com/bob1118/fs/routers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: serverCmdRun,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serverCmd.Flags().BoolP("run", "r", false, "gateway run")
}

func serverCmdRun(cmd *cobra.Command, args []string) {
	fmt.Println("server called")
	if isRun, _ := cmd.Flags().GetBool(`run`); isRun {
		go esl.Run(`inbound`)
		go esl.Run(`outbound`)
		serverHttp()
	} else { //print config gateway
		var list string
		server := viper.Sub(`server`)
		keys := server.AllKeys()
		for _, key := range keys {
			list = fmt.Sprintf("%s\n%-30s=>%s", list, key, server.GetString(key))
		}
		log.Println(`server config:`, list)
	}
}

func serverHttp() {
	db.Initdb()
	h := routers.NewRouter(routers.T_SERVER)
	s := &http.Server{
		Addr:           viper.GetString(`server.http.addr`),
		Handler:        h,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
