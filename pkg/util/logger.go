package util

import (
	"context"
	"github.com/gin-gonic/gin"
	ctxutil "go-server-template/pkg/context"
	"go-server-template/pkg/logger"
)

func Logger(ctx context.Context) logger.Logger {
	var requestID string
	if ginCtx, ok := ctx.(*gin.Context); ok {
		requestID = ctxutil.GetRequestID(ginCtx)
	}
	return logger.GetLogger().WithField("trace_id", requestID)
}
