package core

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-server-template/internal/conf"
	"go-server-template/internal/server/router"
	"go-server-template/pkg/logger"
	"go-server-template/pkg/middleware"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
)

func RunServer() {
	logOut := logger.GetLogger().Writer()
	mux, err := NewMux(
		WithDisablePProf(),
		WithDisableSwagger(),
		WithDisablePrometheus(),
	)

	if err != nil {
		panic(err)
	}

	mws := middleware.New().
		Add("recovery", gin.RecoveryWithWriter(logOut)).
		Add("logger", middleware.LoggerWithConfig(middleware.LoggerConfig{
			Output: logOut,
			// Filter do not add a logger for URLs that contain prefixes such as /debug/, /metrics/, /swagger/, /health
			Filter: func(ctx *gin.Context) bool {
				re, _ := regexp.Compile("^/debug/|^/metrics/|^/swagger/|^/health")
				return re.MatchString(ctx.Request.URL.Path)
			},
		})).All()

	router.Load(mux, mws...)

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
			logger.GetLogger().Warnf("Shutting down due to ServerApi error: %s", err.Error())
		}
	case <-quit:
		break
	}

	logger.GetLogger().Infof("Shutting down server...")

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
