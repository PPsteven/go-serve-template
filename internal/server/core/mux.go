package core

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-server-template/internal/conf"
	"net/http"
	"time"
)

type Option func(option *option)

type option struct {
	disablePProf      bool
	disableSwagger    bool
	disablePrometheus bool
	enableCors        bool
	enableRate        bool
	enableOpenBrowser string
	//alertNotify       proposal.NotifyHandler
	//recordHandler     proposal.RecordHandler
}

// WithDisablePProf 禁用 pprof
func WithDisablePProf() Option {
	return func(opt *option) {
		opt.disablePProf = true
	}
}

// WithDisableSwagger 禁用 swagger
func WithDisableSwagger() Option {
	return func(opt *option) {
		opt.disableSwagger = true
	}
}

// WithDisablePrometheus 禁用prometheus
func WithDisablePrometheus() Option {
	return func(opt *option) {
		opt.disablePrometheus = true
	}
}

// WithAlertNotify 设置告警通知
//func WithAlertNotify(notifyHandler proposal.NotifyHandler) Option {
//	return func(opt *option) {
//		opt.alertNotify = notifyHandler
//	}
//}

// WithRecordMetrics 设置记录接口指标
//func WithRecordMetrics(recordHandler proposal.RecordHandler) Option {
//	return func(opt *option) {
//		opt.recordHandler = recordHandler
//	}
//}

// WithEnableOpenBrowser 启动后在浏览器中打开 uri
func WithEnableOpenBrowser(uri string) Option {
	return func(opt *option) {
		opt.enableOpenBrowser = uri
	}
}

// WithEnableCors 设置支持跨域
func WithEnableCors() Option {
	return func(opt *option) {
		opt.enableCors = true
	}
}

// WithEnableRate 设置支持限流
func WithEnableRate() Option {
	return func(opt *option) {
		opt.enableRate = true
	}
}

func NewOption() *option {
	return &option{}
}

func NewMux(options ...Option) (mux *gin.Engine, err error) {
	mux = gin.New()
	if conf.Conf.Env == conf.Production {
		gin.SetMode(gin.ReleaseMode)
	}
	//mux.StaticFS("assets", http.FS(assets.Bootstrap))
	//mux.SetHTMLTemplate(template.Must(template.New("").ParseFS(assets.Templates, "templates/**/*")))

	opt := new(option)
	for _, f := range options {
		f(opt)
	}

	if !opt.disablePProf {
		if conf.Conf.Env != conf.Production {
			pprof.Register(mux)
		}
	}

	if !opt.disableSwagger {
		if conf.Conf.Env != conf.Production {
			mux.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		}
	}

	mux.NoMethod(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API method.")
	})

	mux.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	// health check
	mux.GET("/health", func(c *gin.Context) {
		resp := &struct {
			Timestamp   time.Time `json:"timestamp"`
			Environment string    `json:"environment"`
			Host        string    `json:"host"`
			Status      string    `json:"status"`
		}{
			Timestamp:   time.Now(),
			Environment: string(conf.Conf.Env),
			Host:        "",
			Status:      "ok",
		}
		c.JSON(http.StatusOK, resp)
	})
	return
}
