package errcode

import (
	"net/http"
)

// Business Code 命名规则
// 10(aa)00(bb)01(cc)
// aa: 通用模块 01-通用模块 02-业务模块
// bb: 业务模块号 01-用户模块
// cc: 具体错误
var (
	ErrParams   = NewSvrError(10001, "params error", http.StatusBadRequest)
	ErrInternal = NewSvrError(10002, "internal error", http.StatusInternalServerError)
)
