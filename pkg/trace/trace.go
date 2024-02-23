package trace

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const Header = "_TRACE_"

var _ T = (*Trace)(nil)

type T interface {
	i()
	ID() string
	WithRequest(*gin.Context) *Trace
	WithResponse(*gin.Context) *Trace
	//AppendDialog(dialog *Dialog) *Trace
	AppendSQL(sql *SQL) *Trace
	//AppendRedis(redis *Redis) *Trace
}

// Trace 记录的参数
type Trace struct {
	mux sync.Mutex
	// TraceID is a unique identifier for the trace.
	Identifier string `json:"trace_id"`
	// Request request info
	Request *Request `json:"request"`
	// Response response info
	Response *Response `json:"response"`
	//ThirdPartyRequests []*Dialog `json:"third_party_requests"` // 调用第三方接口的信息
	//Debugs             []*Debug  `json:"debugs"`               // 调试信息
	SQLs []*SQL `json:"sqls"` // 执行的 SQL 信息
	//Redis              []*Redis  `json:"redis"`                // 执行的 Redis 信息
	// Success shows if the request is successful.
	Success bool `json:"success"`
	// RequestAt show the time of the request.
	RequestAt time.Time `json:"request_at"`
	// ResponseAt show the time of the response.
	ResponseAt time.Time `json:"response_at"`
	// Latency is the duration between the request and the response.
	Latency time.Duration `json:"latency"`
}

// Request 请求信息
type Request struct {
	ClientIP  string      `json:"client_ip"`  // ClientIP equals Context's ClientIP method.
	Method    string      `json:"method"`     // Method is the HTTP method given to the request.
	DecodeURL string      `json:"decode_url"` // DecodeURL is the url after decode
	Header    interface{} `json:"header"`     // 请求 Header 信息
	Body      interface{} `json:"body"`       // 请求 Body 信息
}

// Response 响应信息
type Response struct {
	Header          interface{} `json:"header"`                      // Header 信息
	Body            interface{} `json:"body"`                        // Body 信息
	BusinessCode    int         `json:"business_code,omitempty"`     // 业务码
	BusinessCodeMsg string      `json:"business_code_msg,omitempty"` // 提示信息
	HttpCode        int         `json:"http_code"`                   // HTTP 状态码
	HttpCodeMsg     string      `json:"http_code_msg"`               // HTTP 状态码信息
}

func New(id string) *Trace {
	if id == "" {
		buf := make([]byte, 10)
		io.ReadFull(rand.Reader, buf)
		id = hex.EncodeToString(buf)
	}

	return &Trace{
		Identifier: id,
	}
}

func (t *Trace) i() {}

// ID 唯一标识符
func (t *Trace) ID() string {
	return t.Identifier
}

// WithRequest 设置request
func (t *Trace) WithRequest(c *gin.Context) *Trace {
	decodedURL, _ := url.QueryUnescape(c.Request.URL.RequestURI())

	t.Request = &Request{
		ClientIP:  c.ClientIP(),
		Method:    c.Request.Method,
		DecodeURL: decodedURL,
		Body:      c.Request.Body,
		Header:    c.Request.Header,
	}
	t.RequestAt = time.Now()

	return t
}

// WithResponse 设置response
func (t *Trace) WithResponse(c *gin.Context) *Trace {
	t.Response = &Response{
		Header: c.Writer.Header(),
		//Body:            c.Writer.,
		BusinessCode:    0,
		BusinessCodeMsg: "",
		HttpCode:        c.Writer.Status(),
		HttpCodeMsg:     http.StatusText(c.Writer.Status()),
	}
	t.ResponseAt = time.Now()
	// Calculates the latency.
	t.Latency = t.ResponseAt.Sub(t.RequestAt)

	return t
}

// AppendDialog 安全的追加内部调用过程dialog
//func (t *Trace) AppendDialog(dialog *Dialog) *Trace {
//	if dialog == nil {
//		return t
//	}
//
//	t.mux.Lock()
//	defer t.mux.Unlock()
//
//	t.ThirdPartyRequests = append(t.ThirdPartyRequests, dialog)
//	return t
//}

// AppendDebug 追加 debug
//func (t *Trace) AppendDebug(debug *Debug) *Trace {
//	if debug == nil {
//		return t
//	}
//
//	t.mux.Lock()
//	defer t.mux.Unlock()
//
//	t.Debugs = append(t.Debugs, debug)
//	return t
//}

// AppendSQL 追加 SQL
func (t *Trace) AppendSQL(sql *SQL) *Trace {
	if sql == nil {
		return t
	}

	t.mux.Lock()
	defer t.mux.Unlock()

	t.SQLs = append(t.SQLs, sql)
	return t
}

// AppendRedis 追加 Redis
//func (t *Trace) AppendRedis(redis *Redis) *Trace {
//	if redis == nil {
//		return t
//	}
//
//	t.mux.Lock()
//	defer t.mux.Unlock()
//
//	t.Redis = append(t.Redis, redis)
//	return t
//}

func GetTrace(c *gin.Context) *Trace {
	if v, ok := c.Get(Header); ok {
		if trace, ok := v.(*Trace); ok {
			return trace
		}
	}
	return nil
}
