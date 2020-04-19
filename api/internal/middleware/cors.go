package middleware

import (
	"net/http"
)

type Cors struct {
	handler http.Handler
	Origin  string
}

// ApplyCors applies to Handler the cors config
func ApplyCors(handlerToApply http.Handler, c *Cors) http.Handler {
	c.handler = handlerToApply
	return c
}

func (c *Cors) middleware() *Cors {
	return c
}

func (c *Cors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", c.Origin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	c.handler.ServeHTTP(w, r)
}
