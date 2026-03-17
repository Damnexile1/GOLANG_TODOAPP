package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/logger"
	core_http_middleware "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/middleware"
	core_http_server "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/server"
	user_transport_http "github.com/Damnexile1/GOLANG_TODOAPP/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to initialize logger", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("starting server")
	usersTransportHTTP := user_transport_http.NewUsersHTTPHandler(nil)
	userRoutes := usersTransportHTTP.Routes()
	apiVersionRouter := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(userRoutes...)
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		// XXX: middlewares
		core_http_middleware.RequestId(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)
	httpServer.RegisterApiRoutes(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("Http server run error", zap.Error(err))
	}
}
