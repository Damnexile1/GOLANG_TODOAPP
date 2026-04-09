package core_http_middleware

import (
	"net/http"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
	core_logger "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/logger"
	core_http_response "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/response"
)

// RequireRole middleware проверяет, что у пользователя есть необходимая роль
func RequireRole(requiredRole domain.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := core_logger.FromContext(r.Context())
			responseHandler := core_http_response.NewHTTPResponseHandler(logger, w)

			// Получаем роль из контекста (должна быть установлена Auth middleware)
			userRole, ok := GetUserRoleFromContext(r.Context())
			if !ok {
				responseHandler.ErrorResponse(
					http.ErrNoCookie,
					"user role not found in context",
				)
				return
			}

			role, err := domain.UserRoleFromInt(userRole)
			if err != nil {
				responseHandler.ErrorResponse(
					err,
					"invalid user role",
				)
				return
			}

			// Проверяем права доступа
			if !role.HasPermission(requiredRole) {
				responseHandler.ErrorResponse(
					http.ErrNoCookie,
					"insufficient permissions",
				)
				return
			}

			// Передаем управление следующему handler
			next.ServeHTTP(w, r)
		})
	}
}
