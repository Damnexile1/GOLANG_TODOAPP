package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/auth/jwt"
	core_config "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/config"
	core_logger "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/logger"
	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/middleware"
	core_http_server "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/server"
	auth_service "github.com/Damnexile1/GOLANG_TODOAPP/internal/features/auth/service"
	auth_transport_http "github.com/Damnexile1/GOLANG_TODOAPP/internal/features/auth/transport/http"
	statistics_postgres_repository "github.com/Damnexile1/GOLANG_TODOAPP/internal/features/statistics/repository/postgres"
	statistics_service "github.com/Damnexile1/GOLANG_TODOAPP/internal/features/statistics/service"
	statistics_transport_http "github.com/Damnexile1/GOLANG_TODOAPP/internal/features/statistics/transport/http"
	tasks_postgres_postgres "github.com/Damnexile1/GOLANG_TODOAPP/internal/features/tasks/repository/postgres"
	tasks_service "github.com/Damnexile1/GOLANG_TODOAPP/internal/features/tasks/service"
	tasks_transport_http "github.com/Damnexile1/GOLANG_TODOAPP/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/Damnexile1/GOLANG_TODOAPP/internal/features/users/repository/postgres"
	user_service "github.com/Damnexile1/GOLANG_TODOAPP/internal/features/users/service"
	user_transport_http "github.com/Damnexile1/GOLANG_TODOAPP/internal/features/users/transport/http"
	"go.uber.org/zap"

	_ "github.com/Damnexile1/GOLANG_TODOAPP/docs"
)

// @title    Golang To do API
// @version  1.0
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to initialize logger", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("time_zone", time.Local))

	logger.Debug("initializing postgres connection pool")

	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init connection poo: ", zap.Error(err))
	}

	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := user_service.NewUsersService(usersRepository)
	usersTransportHTTP := user_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing JWT manager")
	jwtConfig := core_config.NewJWTConfigMust()
	jwtManager := jwt.NewJWTManager(jwtConfig.SecretKey, jwtConfig.AccessTokenTTL, jwtConfig.RefreshTokenTTL)

	logger.Debug("initializing feature", zap.String("feature", "auth"))
	authService := auth_service.NewAuthService(usersRepository, jwtManager)
	authTransportHTTP := auth_transport_http.NewAuthHTTPHandler(authService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_postgres.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHttp := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("initializing feature", zap.String("feature", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService)

	logger.Debug("initializing http server")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.CORS(),
		core_http_middleware.RequestId(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(logger),
		core_http_middleware.Trace(),
	)

	// Один роутер для всех endpoints
	// Auth endpoints - публичные (без middleware)
	// Остальные endpoints будут защищены через middleware на уровне роутера
	apiVersionRouterV1 := core_http_server.NewApiVersionRouter(
		core_http_server.ApiVersion1,
		core_http_middleware.Auth(jwtManager), // Применяем Auth ко всем, кроме auth endpoints
	)

	// Регистрируем все routes
	apiVersionRouterV1.RegisterRoutes(authTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(tasksTransportHttp.Routes()...)
	apiVersionRouterV1.RegisterRoutes(statisticsTransportHTTP.Routes()...)

	httpServer.RegisterApiRoutes(apiVersionRouterV1)
	httpServer.RegisterSwagger()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("Http server run error", zap.Error(err))
	}
}
