package core_http_middleware

import (
	"net/http"
	"net/url"
	"runtime/debug"
	"time"

	core_logger "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/logger"
	core_http_response "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const requestIDHeader = "X-Request-ID"

func CORS() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if isAllowedOrigin(origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")

				requestHeaders := r.Header.Get("Access-Control-Request-Headers")
				if requestHeaders != "" {
					w.Header().Set("Access-Control-Allow-Headers", requestHeaders)
				} else {
					w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID, Accept")
				}
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

	u, err := url.Parse(origin)
	if err != nil {
		return false
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	host := u.Hostname()

	switch host {
	case "localhost", "127.0.0.1", "45.131.41.192":
		return true
	default:
		return false
	}
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
				zap.Int("status_code", rw.StatusCode()),
				zap.Duration("time", time.Since(before)),
			)
		})
	}
}
