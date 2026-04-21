package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"time"

	"go.uber.org/zap"
)

var reBotToken = regexp.MustCompile(`(/bot)[^/]+`)

func maskTokenInPath(path string) string {
	return reBotToken.ReplaceAllString(path, "${1}***")
}

// NewTelegramProxy returns a reverse-proxy handler that forwards /tg/{rest}
// to baseURL/{rest}, stripping the /tg prefix.
func NewTelegramProxy(log *zap.Logger, baseURL string) (http.Handler, error) {
	target, err := url.Parse(strings.TrimRight(baseURL, "/"))
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	inner := http.StripPrefix("/tg", proxy)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		inner.ServeHTTP(rw, r)
		log.Info("proxy request",
			zap.String("method", r.Method),
			zap.String("path", maskTokenInPath(r.URL.Path)),
			zap.Int("status", rw.status),
			zap.Duration("duration", time.Since(start)),
		)
	}), nil
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
