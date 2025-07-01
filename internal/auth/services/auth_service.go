package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"voting-blockchain/internal/auth/dto"
	"voting-blockchain/internal/auth/models"
	"voting-blockchain/internal/auth/repositories"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (*models.User, error)
	Login(ctx context.Context, email, password string) (*dto.LoginResponse, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error)
}

type authService struct {
	userRepo         repositories.UserRepository
	refreshRepo      repositories.RefreshTokenRepository
	jwtSecret        string
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
}

func NewAuthService(
	userRepo repositories.UserRepository,
	refreshRepo repositories.RefreshTokenRepository,
	jwtSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) AuthService {
	return &authService{
		userRepo:         userRepo,
		refreshRepo:      refreshRepo,
		jwtSecret:        jwtSecret,
		accessTokenTTL:   accessTTL,
		refreshTokenTTL:  refreshTTL,
	}
}

func (s *authService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *authService) Register(ctx context.Context, email, password string) (*models.User, error) {
	existing, _ := s.userRepo.FindByEmail(ctx, email)
	if existing != nil {
		return nil, errors.New("пользователь с таким email уже существует")
	}

	if len(password) < 8 {
		return nil, errors.New("пароль должен быть не короче 8 символов")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: string(hash),
		Role:         "user", // ← по умолчанию обычный пользователь
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("неверный email или пароль")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("неверный email или пароль")
	}

	accessToken, err := generateJWT(user.ID, user.Role, s.jwtSecret, s.accessTokenTTL)
	if err != nil {
		return nil, err
	}

	refreshToken := generateRandomToken()

	rt := &models.RefreshToken{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(s.refreshTokenTTL),
		Revoked:   false,
	}

	if err := s.refreshRepo.Save(ctx, rt); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	rt, err := s.refreshRepo.FindByToken(ctx, refreshToken)
	if err != nil || rt.Revoked || rt.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token недействителен")
	}

	user, err := s.userRepo.FindByID(ctx, rt.UserID)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	accessToken, err := generateJWT(user.ID, user.Role, s.jwtSecret, s.accessTokenTTL)
	if err != nil {
		return nil, err
	}

	return &dto.RefreshResponse{AccessToken: accessToken}, nil
}

func generateJWT(userID int, role string, secret string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role, // ← добавлено
		"exp":     time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func generateRandomToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
