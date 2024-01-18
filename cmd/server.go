package cmd

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go-server-template/internal/conf"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-server-template/internal/bootstrap"
	"go-server-template/server"
)

var serverCmd = &cobra.Command{
	Use: "serve",
	Short: "start http server with configured api",
	Long: "starts a http server and serves the configured api",
	Run: func(cmd *cobra.Command, args []string) {
		bootstrap.Init()

		serverStart()
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.
	viper.SetDefault("port", "3000")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func serverStart() {
	r := gin.New()
	r.Use(gin.LoggerWithWriter(log.StandardLogger().Out), gin.RecoveryWithWriter(log.StandardLogger().Out))
	server.Init(r)

	var g errgroup.Group
	httpServer := &http.Server{
		Addr: fmt.Sprintf("%s:%d", "localhost", conf.Conf.Port),
		Handler: r,
	}

	g.Go(func() error {
		return httpServer.ListenAndServe()
	})
	serverApiWait := make(chan error)
	go func() {
		serverApiWait <- g.Wait()
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 1 second.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <- serverApiWait:
		if err != nil {
			log.Warnf("Shutting down due to ServerApi error: %s", err.Error())
		}
	case <- quit:
		break
	}

	log.Infof("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := httpServer.Shutdown(ctx)
	if err == nil {
		log.Println("Server gracefully stopped")
	} else {
		log.Fatalf("Failed to shutdown http server: %v", err)
	}
	return
}