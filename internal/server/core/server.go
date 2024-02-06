package core

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go-server-template/internal/conf"
	"go-server-template/internal/server/router"
	"go-server-template/pkg/middleware"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func RunServer() {
	logOut := log.StandardLogger().Out
	mux, err := NewMux(
		WithCustomMiddleware("recovery", gin.RecoveryWithWriter(logOut)),
		WithCustomMiddleware("logger", middleware.LoggerWithConfig(middleware.LoggerConfig{
			Output: logOut,
			SkipPaths: []string{
				"/debug/pprof/",
				"/debug/pprof/cmdline",
				"/debug/pprof/profile",
				"/debug/pprof/symbol",
				"/debug/pprof/symbol",
				"/debug/pprof/trace",
				"/debug/pprof/allocs",
				"/debug/pprof/block",
				"/debug/pprof/goroutine",
				"/debug/pprof/heap",
				"/debug/pprof/mutex",
				"/debug/pprof/threadcreate",
			},
		})),
	)

	if err != nil {
		panic(err)
	}

	router.Init(mux)

	var g errgroup.Group
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", "localhost", conf.Conf.Port),
		Handler: mux,
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
	case err := <-serverApiWait:
		if err != nil {
			log.Warnf("Shutting down due to ServerApi error: %s", err.Error())
		}
	case <-quit:
		break
	}

	log.Infof("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = httpServer.Shutdown(ctx)
	if err == nil {
		log.Println("Server gracefully stopped")
	} else {
		log.Fatalf("Failed to shutdown http server: %v", err)
	}
	return
}
