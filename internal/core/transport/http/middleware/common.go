package core_http_middleware

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	core_logger "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/logger"
	core_http_response "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	requestIDHeader = "X-Request-ID"
)

func CORS() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
			if err != nil {
				fmt.Println("failed to initialize logger", err)
				os.Exit(1)
			}
			defer logger.Close()

			logger.Debug("cors check", zap.String("origin", origin), zap.String("method", r.Method))
			if isAllowedOrigin(origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isAllowedOrigin(origin string) bool {
	if origin == "" {
		return false
	}

	if strings.HasPrefix(origin, "http://localhost:") {
		return true
	}

	if strings.HasPrefix(origin, "http://127.0.0.1:") {
		return true
	}

	if strings.HasPrefix(origin, "http://45.131.41.192/:") {
		return true
	}

	return false
}

func RequestId() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := r.Header.Get(requestIDHeader)
			if requestId == "" {
				requestId = uuid.NewString()
			}
			r.Header.Set(requestIDHeader, requestId)
			w.Header().Set(requestIDHeader, requestId)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := r.Header.Get(requestIDHeader)
			l := log.With(
				zap.String("requestId", requestId),
				zap.String("url", r.URL.String()),
			)

			ctx := core_logger.ContextWithLogger(r.Context(), l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Panic(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					log.Error(
						"panic recovered",
						zap.Any("panic", p),
						zap.ByteString("stack", debug.Stack()),
					)

					responseHandler.PanicResponse(
						p,
						"during handle HTTP request got unexpected panic",
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			rw := core_http_response.NewResponseWriter(w)

			before := time.Now()
			log.Debug(
				">>> handling HTTP request",
				zap.String("method", r.Method),
				zap.Time("time", before.UTC()),
			)

			next.ServeHTTP(rw, r.WithContext(ctx))

			log.Debug(
				"<<< finished HTTP request",
				zap.Int("status_code: ", rw.StatusCode()),
				zap.Duration("time", time.Now().Sub(before)),
			)
		})
	}
}
