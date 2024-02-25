package errcode

import (
	"net/http"
)

var (
	ErrUserNotFound = NewSvrError(200101, "user not found", http.StatusNotFound)
)
