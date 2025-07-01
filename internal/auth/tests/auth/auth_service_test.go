package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"voting-blockchain/internal/auth/models"
	"voting-blockchain/internal/auth/services"

	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	users map[string]*models.User
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	if u, ok := m.users[email]; ok {
		return u, nil
	}
	return nil, nil
}

func (m *mockUserRepo) FindByID(ctx context.Context, id int) (*models.User, error) {
    for _, u := range m.users {
        if u.ID == id {
            return u, nil
        }
    }
    return nil, errors.New("user not found")
}

func mustHashPassword(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		panic("не удалось захешировать пароль в тесте")
	}
	return string(hash)
}


func (m *mockUserRepo) Create(ctx context.Context, u *models.User) error {
	m.users[u.Email] = u
	u.ID = 1
	return nil
}

type mockRefreshRepo struct {
	tokens map[string]*models.RefreshToken
}

func (m *mockRefreshRepo) Save(ctx context.Context, t *models.RefreshToken) error {
	m.tokens[t.Token] = t
	return nil
}
func (m *mockRefreshRepo) FindByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	if t, ok := m.tokens[token]; ok {
		return t, nil
	}
	return nil, nil
}
func (m *mockRefreshRepo) Revoke(ctx context.Context, token string) error {
	if t, ok := m.tokens[token]; ok {
		t.Revoked = true
	}
	return nil
}

func TestRegister(t *testing.T) {
	userRepo := &mockUserRepo{users: map[string]*models.User{}}
	refreshRepo := &mockRefreshRepo{tokens: map[string]*models.RefreshToken{}}

	svc := services.NewAuthService(userRepo, refreshRepo, "secret", 15*time.Minute, 7*24*time.Hour)

	user, err := svc.Register(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatal(err)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %s", user.Email)
	}
}

func TestLogin_Success(t *testing.T) {
	userRepo := &mockUserRepo{users: map[string]*models.User{}}
	refreshRepo := &mockRefreshRepo{tokens: map[string]*models.RefreshToken{}}

	// Регистрация пользователя вручную
	hashed, _ := bcrypt.GenerateFromPassword([]byte("pass1234"), bcrypt.DefaultCost)
	userRepo.users["me@example.com"] = &models.User{
		ID:           1,
		Email:        "me@example.com",
		PasswordHash: string(hashed),
	}

	svc := services.NewAuthService(userRepo, refreshRepo, "secret", 15*time.Minute, 7*24*time.Hour)

	resp, err := svc.Login(context.Background(), "me@example.com", "pass1234")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if resp.AccessToken == "" || resp.RefreshToken == "" {
		t.Error("tokens should not be empty")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	userRepo := &mockUserRepo{users: map[string]*models.User{}}
	refreshRepo := &mockRefreshRepo{tokens: map[string]*models.RefreshToken{}}

	hashed, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
	userRepo.users["me@example.com"] = &models.User{
		ID:           1,
		Email:        "me@example.com",
		PasswordHash: string(hashed),
	}

	svc := services.NewAuthService(userRepo, refreshRepo, "secret", 15*time.Minute, 7*24*time.Hour)

	_, err := svc.Login(context.Background(), "me@example.com", "wrongpass")
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
}

func TestRefresh_Success(t *testing.T) {
	userRepo := &mockUserRepo{}
	refreshRepo := &mockRefreshRepo{tokens: map[string]*models.RefreshToken{}}

	token := "valid-refresh-token"
	refreshRepo.tokens[token] = &models.RefreshToken{
		Token:     token,
		UserID:    1,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Revoked:   false,
	}

	svc := services.NewAuthService(userRepo, refreshRepo, "secret", 15*time.Minute, 7*24*time.Hour)

	resp, err := svc.Refresh(context.Background(), token)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if resp.AccessToken == "" {
		t.Error("expected access token")
	}
}

func TestRefresh_RevokedToken(t *testing.T) {
	userRepo := &mockUserRepo{}
	refreshRepo := &mockRefreshRepo{tokens: map[string]*models.RefreshToken{}}

	token := "revoked-token"
	refreshRepo.tokens[token] = &models.RefreshToken{
		Token:     token,
		UserID:    1,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Revoked:   true,
	}

	svc := services.NewAuthService(userRepo, refreshRepo, "secret", 15*time.Minute, 7*24*time.Hour)

	_, err := svc.Refresh(context.Background(), token)
	if err == nil {
		t.Fatal("expected error for revoked token")
	}
}

func TestRefresh_ExpiredToken(t *testing.T) {
	userRepo := &mockUserRepo{}
	refreshRepo := &mockRefreshRepo{tokens: map[string]*models.RefreshToken{}}

	token := "expired-token"
	refreshRepo.tokens[token] = &models.RefreshToken{
		Token:     token,
		UserID:    1,
		ExpiresAt: time.Now().Add(-1 * time.Hour), // истёк
		Revoked:   false,
	}

	svc := services.NewAuthService(userRepo, refreshRepo, "secret", 15*time.Minute, 7*24*time.Hour)

	_, err := svc.Refresh(context.Background(), token)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}

func TestGetUserByID_Success(t *testing.T) {
    userRepo := &mockUserRepo{
        users: map[string]*models.User{
            "user@example.com": {ID: 1, Email: "user@example.com"},
        },
    }
    refreshRepo := &mockRefreshRepo{}
    svc := services.NewAuthService(userRepo, refreshRepo, "secret", 15*time.Minute, 7*24*time.Hour)

    user, err := svc.GetUserByID(context.Background(), 1)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if user.Email != "user@example.com" {
        t.Errorf("expected email 'user@example.com', got '%s'", user.Email)
    }
}
func TestGetUserByID_NotFound(t *testing.T) {
    userRepo := &mockUserRepo{
        users: map[string]*models.User{},
    }
    refreshRepo := &mockRefreshRepo{}
    svc := services.NewAuthService(userRepo, refreshRepo, "secret", 15*time.Minute, 7*24*time.Hour)

    _, err := svc.GetUserByID(context.Background(), 42)
    if err == nil {
        t.Fatal("expected error for non-existent user")
    }
}

func TestLogin_Success_Variant(t *testing.T) {
    mockUsers := map[string]*models.User{
        "user@example.com": {
            ID:           1,
            Email:        "user@example.com",
            PasswordHash: mustHashPassword("securepassword"),
            CreatedAt:    time.Now(),
        },
    }

    userRepo := &mockUserRepo{users: mockUsers}
    refreshRepo := &mockRefreshRepo{tokens: make(map[string]*models.RefreshToken)}
    service := services.NewAuthService(userRepo, refreshRepo, "secret", 15*time.Minute, 7*24*time.Hour)

    resp, err := service.Login(context.Background(), "user@example.com", "securepassword")
    if err != nil {
        t.Fatalf("ожидался успешный логин, но ошибка: %v", err)
    }

    if resp.AccessToken == "" || resp.RefreshToken == "" {
        t.Fatal("токены не сгенерированы")
    }
}

func TestLogin_InvalidPassword(t *testing.T) {
    mockUsers := map[string]*models.User{
        "user@example.com": {
            ID:           1,
            Email:        "user@example.com",
            PasswordHash: mustHashPassword("securepassword"),
        },
    }

    userRepo := &mockUserRepo{users: mockUsers}
    refreshRepo := &mockRefreshRepo{tokens: make(map[string]*models.RefreshToken)}
    service := services.NewAuthService(userRepo, refreshRepo, "secret", 15*time.Minute, 7*24*time.Hour)

    _, err := service.Login(context.Background(), "user@example.com", "wrongpassword")
    if err == nil {
        t.Fatal("ожидалась ошибка при неправильном пароле")
    }
}
func TestRefresh_TokenNotFound(t *testing.T) {
	refreshRepo := &mockRefreshRepo{
		tokens: map[string]*models.RefreshToken{},
	}
	userRepo := &mockUserRepo{}

	service := services.NewAuthService(userRepo, refreshRepo, "secret", 15*time.Minute, 7*24*time.Hour)

	_, err := service.Refresh(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("ожидалась ошибка при несуществующем токене")
	}
}


