package core_http_middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/auth/jwt"
	core_logger "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/logger"
	core_http_response "github.com/Damnexile1/GOLANG_TODOAPP/internal/core/transport/http/response"
)

type contextKey string

const (
	UserIDContextKey    contextKey = "user_id"
	UserEmailContextKey contextKey = "user_email"
	UserRoleContextKey  contextKey = "user_role"
)

// Auth middleware проверяет JWT токен и добавляет данные пользователя в контекст
// Пропускает публичные пути (auth endpoints)
func Auth(jwtManager *jwt.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Пропускаем auth endpoints без проверки токена
			if strings.HasPrefix(r.URL.Path, "/auth/") {
				next.ServeHTTP(w, r)
				return
			}

			logger := core_logger.FromContext(r.Context())
			responseHandler := core_http_response.NewHTTPResponseHandler(logger, w)

			// Получаем токен из заголовка Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				responseHandler.ErrorResponse(
					http.ErrNoCookie,
					"missing authorization header",
				)
				return
			}

			// Проверяем формат "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				responseHandler.ErrorResponse(
					http.ErrNoCookie,
					"invalid authorization header format",
				)
				return
			}

			tokenString := parts[1]

			// Валидируем токен
			claims, err := jwtManager.ValidateToken(tokenString)
			if err != nil {
				responseHandler.ErrorResponse(
					err,
					"invalid or expired token",
				)
				return
			}

			// Добавляем данные пользователя в контекст
			ctx := r.Context()
			ctx = context.WithValue(ctx, UserIDContextKey, claims.UserID)
			ctx = context.WithValue(ctx, UserEmailContextKey, claims.Email)
			ctx = context.WithValue(ctx, UserRoleContextKey, claims.Role)

			// Передаем управление следующему handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext извлекает user_id из контекста
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(int)
	return userID, ok
}

// GetUserEmailFromContext извлекает email из контекста
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(UserEmailContextKey).(string)
	return email, ok
}

// GetUserRoleFromContext извлекает role из контекста
func GetUserRoleFromContext(ctx context.Context) (int, bool) {
	role, ok := ctx.Value(UserRoleContextKey).(int)
	return role, ok
}
