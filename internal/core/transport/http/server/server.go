package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Damnexile1/GOLANG_TODOAPP/docs"
	core_logger "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/logger"
	core_http_middleware "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/middleware"
	"github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux    *http.ServeMux
	config Config
	log    core_logger.Logger

	middleware []core_http_middleware.Middleware
}

func NewHTTPServer(
	config Config,
	log *core_logger.Logger,
	middleware ...core_http_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        *log,
		middleware: middleware,
	}
}

func (s *HTTPServer) RegisterApiRoutes(routers ...*ApiVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)
		s.log.Debug("registering route", zap.String("prefix", prefix))
		s.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router.WithMiddleware()),
		)
	}
}

func (s *HTTPServer) RegisterSwagger() {
	s.mux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		),
	)
	s.mux.HandleFunc(
		"/swagger/doc.json",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(docs.SwaggerInfo.ReadDoc()))
		},
	)
}

func (s *HTTPServer) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddleware(s.mux, s.middleware...)

	server := &http.Server{
		Addr:    s.config.Address,
		Handler: mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Warn("Listening on ", zap.String("addr", s.config.Address))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("http server: %w", err)
		}
	case <-ctx.Done():
		s.log.Warn("Shutting down HTTP server")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.config.ShutDownTimeout,
		)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("http server shutdown: %w", err)
		}

		s.log.Warn("HTTP server stopped")
		return nil
	}
	return nil
}
