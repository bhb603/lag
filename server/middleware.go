package server

import (
	"fmt"
	"net/http"
	"time"
)

func parseParamsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		next.ServeHTTP(w, r)
	})
}

type lagMiddleware struct {
	maxLag time.Duration
}

func (lm *lagMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if t := r.Form.Get("lag"); len(t) > 0 {
			d, err := time.ParseDuration(t)
			if err != nil {
				http.Error(w, "invalid time parameter", http.StatusBadRequest)
				return
			}
			if d > lm.maxLag {
				http.Error(w, fmt.Sprintf("lag time cannot exceed %s", lm.maxLag), http.StatusBadRequest)
				return
			}
			time.Sleep(d)
		}
		next.ServeHTTP(w, r)
	})

}
