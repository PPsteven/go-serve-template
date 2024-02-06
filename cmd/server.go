package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go-server-template/internal/bootstrap"
	"go-server-template/internal/server/core"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start http server with configured api",
	Long:  "starts a http server and serves the configured api",
	Run: func(cmd *cobra.Command, args []string) {
		bootstrap.Init()

		core.RunServer()
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.
	viper.SetDefault("port", 3000)
	viper.SetDefault("env", "dev")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
