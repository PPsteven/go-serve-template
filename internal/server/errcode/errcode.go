package errcode

import "fmt"

var _ SvrError = (*svrError)(nil)

type SvrError interface {
	Error() string

	Code() int

	Message() string

	HttpCode() int

	Detail() []string

	WithDetail(string, ...interface{}) SvrError

	WithError(err error) SvrError

	i()
}

type svrError struct {
	code     int
	message  string
	httpCode int
	detail   []string
}

func (e *svrError) Error() string {
	return fmt.Sprint("code: ", e.code, " message: ", e.message, " detail: ", e.detail)
}

func (e *svrError) Code() int {
	return e.code
}

func (e *svrError) Message() string {
	return e.message
}

func (e *svrError) HttpCode() int {
	return e.httpCode
}

func (e *svrError) Detail() []string {
	return e.detail
}

func (e *svrError) WithDetail(format string, a ...interface{}) SvrError {
	c := *e
	c.detail = append(c.detail, fmt.Sprintf(format, a))
	return &c
}

func (e *svrError) WithError(err error) SvrError {
	c := *e
	c.detail = append(c.detail, err.Error())
	return &c
}

func (e *svrError) i() {}

func NewSvrError(code int, message string, httpCode int) SvrError {
	return &svrError{
		code:     code,
		message:  message,
		httpCode: httpCode,
	}
}
