package core

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go-server-template/internal/conf"
	"go-server-template/pkg/middleware"
)

type Option func(option *option)

type option struct {
	disablePProf      bool
	disableSwagger    bool
	disablePrometheus bool
	enableCors        bool
	enableRate        bool
	enableOpenBrowser string
	customMiddlewares map[string]gin.HandlerFunc
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

func WithCustomMiddleware(key string, mw gin.HandlerFunc) Option {
	return func(opt *option) {
		opt.customMiddlewares[key] = mw
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
	opt.customMiddlewares = make(map[string]gin.HandlerFunc)
	for _, f := range options {
		f(opt)
	}

	if !opt.disablePProf {
		if conf.Conf.Env != conf.Production {
			pprof.Register(mux)
		}
	}

	for key, mw := range opt.customMiddlewares {
		middleware.Middlewares[key] = mw
	}
	for _, mw := range middleware.Middlewares {
		mux.Use(mw)
	}
	return
}
