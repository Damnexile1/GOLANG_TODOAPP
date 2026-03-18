package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/logger"
	core_postgres_pool "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/middleware"
	core_http_server "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/server"
	users_postgres_repository "github.com/Damnexile1/GOLANG_TODOAPP/internal/features/users/repository/postgres"
	user_service "github.com/Damnexile1/GOLANG_TODOAPP/internal/features/users/service"
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

	logger.Debug("initializing postgres connection pool")
	pool, err := core_postgres_pool.NewConnectionPool(
		ctx,
		core_postgres_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init connection poo: ", zap.Error(err))
	}

	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := user_service.NewUsersService(usersRepository)
	usersTransportHTTP := user_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing http server")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestId(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(logger),
		core_http_middleware.Trace(),
	)
	apiVersionRouter := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	httpServer.RegisterApiRoutes(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("Http server run error", zap.Error(err))
	}
}
