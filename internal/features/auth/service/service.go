package auth_service

import (
	"context"
	"fmt"

	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/auth/jwt"
	"github.com/Damnexile1/GOLANG_TODOAPP/internal/core/domain"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	usersRepository UsersRepository
	jwtManager      *jwt.JWTManager
}

type UsersRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}

func NewAuthService(usersRepository UsersRepository, jwtManager *jwt.JWTManager) *AuthService {
	return &AuthService{
		usersRepository: usersRepository,
		jwtManager:      jwtManager,
	}
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(ctx context.Context, email, password, fullName string, phoneNumber *string) (*jwt.TokenPair, error) {
	// Проверяем, не существует ли уже пользователь с таким email
	_, err := s.usersRepository.GetUserByEmail(ctx, email)
	if err == nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Хешируем пароль
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	// Создаем пользователя
	user := domain.NewUserUninitialized(
		fullName,
		phoneNumber,
		email,
		string(passwordHash),
		domain.UserRoleUser, // По умолчанию обычный пользователь
		nil,
	)

	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("validate user: %w", err)
	}

	user, err = s.usersRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	// Генерируем токены
	tokens, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return tokens, nil
}

// Login авторизует пользователя
func (s *AuthService) Login(ctx context.Context, email, password string) (*jwt.TokenPair, error) {
	// Получаем пользователя по email
	user, err := s.usersRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Генерируем токены
	tokens, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return tokens, nil
}

// RefreshToken обновляет access token используя refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*jwt.TokenPair, error) {
	// Валидируем refresh token
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Получаем актуальные данные пользователя
	user, err := s.usersRepository.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Генерируем новую пару токенов
	tokens, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return tokens, nil
}
