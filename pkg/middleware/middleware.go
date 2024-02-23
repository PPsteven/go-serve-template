package middleware

import "github.com/gin-gonic/gin"

type Middleware struct {
	middlewares []map[string]gin.HandlerFunc
}

func defaultMiddlewares() []map[string]gin.HandlerFunc {
	return []map[string]gin.HandlerFunc{
		{"recovery": gin.Recovery()},
		{"secure": Secure},
		{"options": Options},
		{"nocache": NoCache},
		{"request_id": RequestID()},
		{"logger": Logger()},
		//"trace":      Tracing,
		{"cors": Cors()},
	}
}

func New() *Middleware {
	return &Middleware{middlewares: defaultMiddlewares()}
}

func (m *Middleware) Add(name string, handleFunc gin.HandlerFunc) *Middleware {
	for i := range m.middlewares {
		if _, ok := m.middlewares[i][name]; ok {
			m.middlewares[i][name] = handleFunc
			return m
		}
	}

	m.middlewares = append(m.middlewares, map[string]gin.HandlerFunc{name: handleFunc})
	return m
}

func (m *Middleware) All() []gin.HandlerFunc {
	mws := make([]gin.HandlerFunc, len(m.middlewares))
	for i, mw := range m.middlewares {
		mws[i] = getHandler(mw)
	}
	return mws
}

func getHandler(handler map[string]gin.HandlerFunc) gin.HandlerFunc {
	for _, v := range handler {
		return v
	}
	return nil
}
