package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/unrolled/slog"
)

func basicAuth(handler http.HandlerFunc, allowedUser, allowedPassword string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		user, password, ok := req.BasicAuth()
		if !ok || len(strings.TrimSpace(user)) == 0 || len(strings.TrimSpace(password)) == 0 {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if user != allowedUser || password != allowedPassword {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		handler(w, req)
	}
}

func noCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "max-age=0, private, must-revalidate")
		next.ServeHTTP(w, r)
	})
}

func recovery(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

				stack := make([]byte, 8*1024)
				stack = stack[:runtime.Stack(stack, false)]
				slog.Error("panic_recovery", slog.String("err", fmt.Sprintf("%v", err)), slog.String("stack", string(stack)))
			}
		}()

		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}

func log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		crw := newCustomResponseWriter(w)
		next.ServeHTTP(crw, r)

		end := time.Since(start)

		slog.Info("request",
			slog.Request(r),
			slog.String("address", r.RemoteAddr),
			slog.String("method", r.Method),
			slog.String("uri", r.RequestURI),
			slog.Int("status", crw.status),
			slog.Int("size", crw.size),
			slog.NullableString("token", r.Header.Get("X-DeviceToken")),
			slog.Duration("duration", end),
			slog.String("human", end.String()),
		)
	})
}

type customResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (c *customResponseWriter) WriteHeader(status int) {
	c.status = status
	c.ResponseWriter.WriteHeader(status)
}

func (c *customResponseWriter) Write(b []byte) (int, error) {
	size, err := c.ResponseWriter.Write(b)
	c.size += size
	return size, err
}

func newCustomResponseWriter(w http.ResponseWriter) *customResponseWriter {
	return &customResponseWriter{
		ResponseWriter: w,
		status:         200,
	}
}
