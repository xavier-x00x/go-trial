package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go-trial/internal/config"
	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/domain/entity"
	"go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"
	jwtPkg "go-trial/pkg/jwt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	ErrEmailAlreadyExists    = errors.New("email already registered")
	ErrUsernameAlreadyExists = errors.New("username already taken")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrInvalidRefreshToken   = errors.New("invalid or expired refresh token")
	ErrUserNotFound          = errors.New("user not found")
	ErrAccountInactive       = errors.New("account is inactive")
	ErrRoleNotAssigned       = errors.New("account is pending approval: role not assigned")
)

// AuthUseCase defines the business logic for authentication.
type AuthUseCase interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, string, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, string, error)
	GoogleLogin(ctx context.Context, code string) (*dto.AuthResponse, string, error)
	GetGoogleLoginURL() string
	RefreshToken(ctx context.Context, refreshTokenStr string) (*dto.RefreshResponse, error)
	GetMe(ctx context.Context, userID string) (*dto.UserResponse, error)
	GetAllUsers(ctx context.Context) ([]dto.UserResponse, error)
	GetAllUsersWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.UserResponse, *entity.Meta, error)
	UpdateUser(ctx context.Context, id string, req dto.UpdateUserByAdminRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id string) error
}

type authUseCase struct {
	userRepo   repository.UserRepository
	uow        uow.UnitOfWork
	jwtManager *jwtPkg.JWTManager
	googleCfg  *oauth2.Config
}

func NewAuthUseCase(
	userRepo repository.UserRepository,
	uow uow.UnitOfWork,
	jwtManager *jwtPkg.JWTManager,
	cfg *config.Config,
) AuthUseCase {
	googleCfg := &oauth2.Config{
		ClientID:     cfg.GoogleOAuth.ClientID,
		ClientSecret: cfg.GoogleOAuth.ClientSecret,
		RedirectURL:  cfg.GoogleOAuth.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &authUseCase{
		userRepo:   userRepo,
		uow:        uow,
		jwtManager: jwtManager,
		googleCfg:  googleCfg,
	}
}

// Register creates a new user within a UoW transaction and returns tokens.
func (u *authUseCase) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, string, error) {
	// Check if email already exists
	existing, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", err
	}
	if existing != nil {
		return nil, "", ErrEmailAlreadyExists
	}

	// Check if username already exists
	existing, err = u.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, "", err
	}
	if existing != nil {
		return nil, "", ErrUsernameAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// Generate UUID v7
	id, err := uuid.NewV7()
	if err != nil {
		return nil, "", err
	}

	// Default role
	role := req.Role
	if role == "" {
		role = "staff"
	}

	isActive := true
	user := &entity.User{
		ID:       id.String(),
		StoreID:  req.StoreID,
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: string(hashedPassword),
		Role:     role,
		IsActive: &isActive,
	}

	// Begin transaction via UoW
	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, "", err
	}
	defer u.uow.Rollback(txCtx) //nolint:errcheck

	if err := u.userRepo.Create(txCtx, user); err != nil {
		return nil, "", err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, "", err
	}

	// Generate tokens
	accessToken, err := u.jwtManager.GenerateAccessToken(user.ID, user.Email, user.Role, user.StoreID)
	if err != nil {
		return nil, "", err
	}

	refreshToken, err := u.jwtManager.GenerateRefreshToken(user.ID, user.Email, user.Role, user.StoreID)
	if err != nil {
		return nil, "", err
	}

	return &dto.AuthResponse{
		AccessToken: accessToken,
		User:        toUserResponse(user),
	}, refreshToken, nil
}

// Login authenticates a user by email, username, or phone and returns tokens.
func (u *authUseCase) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, string, error) {
	identity := strings.TrimSpace(req.Identity)

	// Try to find user by email, username, or phone
	user, err := u.userRepo.FindByIdentity(ctx, identity)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", ErrInvalidCredentials
	}

	// Check active status
	if user.IsActive != nil && !*user.IsActive {
		return nil, "", ErrAccountInactive
	}

	// Check if role is assigned
	if user.Role == "" {
		return nil, "", ErrRoleNotAssigned
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	_ = u.userRepo.Update(ctx, user)

	accessToken, err := u.jwtManager.GenerateAccessToken(user.ID, user.Email, user.Role, user.StoreID)
	if err != nil {
		return nil, "", err
	}

	var refreshToken string
	if req.Remember != nil && *req.Remember {
		refreshToken, err = u.jwtManager.GenerateRememberMeToken(user.ID, user.Email, user.Role, user.StoreID)
	} else {
		refreshToken, err = u.jwtManager.GenerateRefreshToken(user.ID, user.Email, user.Role, user.StoreID)
	}
	if err != nil {
		return nil, "", err
	}

	return &dto.AuthResponse{
		AccessToken: accessToken,
		User:        toUserResponse(user),
	}, refreshToken, nil
}

// RefreshToken validates the refresh token and issues a new access token.
func (u *authUseCase) RefreshToken(ctx context.Context, refreshTokenStr string) (*dto.RefreshResponse, error) {
	claims, err := u.jwtManager.ValidateToken(refreshTokenStr)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	if claims.Type != jwtPkg.RefreshToken {
		return nil, ErrInvalidRefreshToken
	}

	// Verify user still exists and is active
	user, err := u.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if user.IsActive != nil && !*user.IsActive {
		return nil, ErrAccountInactive
	}

	accessToken, err := u.jwtManager.GenerateAccessToken(user.ID, user.Email, user.Role, user.StoreID)
	if err != nil {
		return nil, err
	}

	return &dto.RefreshResponse{
		AccessToken: accessToken,
	}, nil
}

// GetMe returns the current user's profile.
func (u *authUseCase) GetMe(ctx context.Context, userID string) (*dto.UserResponse, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	resp := toUserResponse(user)
	return &resp, nil
}

func (u *authUseCase) GetAllUsers(ctx context.Context) ([]dto.UserResponse, error) {
	users, err := u.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var resp []dto.UserResponse
	for _, user := range users {
		resp = append(resp, toUserResponse(&user))
	}
	return resp, nil
}

func (u *authUseCase) GetAllUsersWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.UserResponse, *entity.Meta, error) {
	allowedOrder := []string{"id", "name", "username", "email", "created_at"}
	searchColumns := []string{"id", "name", "username", "email"}

	filter := BuildQueryFilter(meta, allowedOrder, searchColumns)

	data, resMeta, err := u.userRepo.FindAllWithPagination(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	var resp []dto.UserResponse
	for _, user := range data {
		resp = append(resp, toUserResponse(&user))
	}
	return resp, resMeta, nil
}

func (u *authUseCase) UpdateUser(ctx context.Context, id string, req dto.UpdateUserByAdminRequest) (*dto.UserResponse, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Username != nil {
		existing, _ := u.userRepo.FindByUsername(ctx, *req.Username)
		if existing != nil && existing.ID != user.ID {
			return nil, ErrUsernameAlreadyExists
		}
		user.Username = *req.Username
	}
	if req.Email != nil {
		existing, _ := u.userRepo.FindByEmail(ctx, *req.Email)
		if existing != nil && existing.ID != user.ID {
			return nil, ErrEmailAlreadyExists
		}
		user.Email = *req.Email
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.StoreID != nil {
		user.StoreID = req.StoreID
	}
	if req.IsActive != nil {
		user.IsActive = req.IsActive
	}
	if req.Password != nil && *req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.userRepo.Update(txCtx, user); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	resp := toUserResponse(user)
	return &resp, nil
}

func (u *authUseCase) DeleteUser(ctx context.Context, id string) error {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.userRepo.Delete(txCtx, id); err != nil {
		return err
	}

	return u.uow.Commit(txCtx)
}

func toUserResponse(user *entity.User) dto.UserResponse {
	return dto.UserResponse{
		ID:          user.ID,
		StoreID:     user.StoreID,
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		Phone:       user.Phone,
		Role:        user.Role,
		AvatarURL:   user.AvatarURL,
		IsActive:    user.IsActive,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

type googleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func (u *authUseCase) GoogleLogin(ctx context.Context, code string) (*dto.AuthResponse, string, error) {
	// 1. Exchange code for token
	token, err := u.googleCfg.Exchange(ctx, code)
	if err != nil {
		return nil, "", fmt.Errorf("failed to exchange code: %w", err)
	}

	// 2. Fetch user info from Google
	client := u.googleCfg.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var gUser googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&gUser); err != nil {
		return nil, "", fmt.Errorf("failed to decode user info: %w", err)
	}

	// 3. Find or Create user
	// Try by GoogleID first
	user, err := u.userRepo.FindByGoogleID(ctx, gUser.ID)
	if err != nil {
		return nil, "", err
	}

	// If not found by GoogleID, try by Email (Link account)
	if user == nil {
		user, err = u.userRepo.FindByEmail(ctx, gUser.Email)
		if err != nil {
			return nil, "", err
		}
	}

	isNewUser := false
	if user == nil {
		// Create new user
		isNewUser = true
		id, err := uuid.NewV7()
		if err != nil {
			return nil, "", err
		}

		user = &entity.User{
			ID:           id.String(),
			Name:         gUser.Name,
			Email:        gUser.Email,
			Username:     strings.Split(gUser.Email, "@")[0], // Fallback username
			Password:     "password",                         // Default password for google users
			Role:         "",                                 // New google users have NO role
			AuthProvider: "google",
			GoogleID:     &gUser.ID,
		}
		// Since username must be unique, check if exists and append something if needed
		existingUser, _ := u.userRepo.FindByUsername(ctx, user.Username)
		if existingUser != nil {
			user.Username = fmt.Sprintf("%s_%s", user.Username, id.String()[:8])
		}
	} else {
		// Link GoogleID if not linked
		if user.GoogleID == nil {
			user.GoogleID = &gUser.ID
			user.AuthProvider = "google"
		}
	}

	// Check if role is assigned (Requirement: cannot login if role is empty)
	if user.Role == "" {
		// Even if new user, we still save it so Admin can see it in user list
		// BUT we must save it first before returning error, or at least have a way to see it.
		// Actually, let's save the progress so far (download picture, update last login etc)
		// but skip token generation if no role.
	}

	// Download/Sync Profile Picture
	if gUser.Picture != "" {
		localAvatarPath, err := u.downloadGoogleProfilePicture(user.ID, gUser.Picture)
		if err == nil {
			user.AvatarURL = &localAvatarPath
		}
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now

	// Begin Transaction via UoW
	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, "", err
	}
	defer u.uow.Rollback(txCtx)

	if isNewUser {
		if err := u.userRepo.Create(txCtx, user); err != nil {
			return nil, "", err
		}
	} else {
		if err := u.userRepo.Update(txCtx, user); err != nil {
			return nil, "", err
		}
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, "", err
	}

	// CHECK ROLE AGAIN AFTER COMMIT (to ensure we returned the latest state from DB if needed)
	if user.Role == "" {
		return nil, "", ErrRoleNotAssigned
	}

	// 4. Generate Tokens
	accessToken, err := u.jwtManager.GenerateAccessToken(user.ID, user.Email, user.Role, user.StoreID)
	if err != nil {
		return nil, "", err
	}

	refreshToken, err := u.jwtManager.GenerateRefreshToken(user.ID, user.Email, user.Role, user.StoreID)
	if err != nil {
		return nil, "", err
	}

	return &dto.AuthResponse{
		AccessToken: accessToken,
		User:        toUserResponse(user),
	}, refreshToken, nil
}

func (u *authUseCase) GetGoogleLoginURL() string {
	return u.googleCfg.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

func (u *authUseCase) downloadGoogleProfilePicture(userID, pictureURL string) (string, error) {
	if pictureURL == "" {
		return "", nil
	}

	resp, err := http.Get(pictureURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download picture: %s", resp.Status)
	}

	// Create directory if not exists
	dir := "uploads/avatars"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	// Generate filename
	ext := ".jpg" // default
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "png") {
		ext = ".png"
	} else if strings.Contains(contentType, "webp") {
		ext = ".webp"
	}

	filename := fmt.Sprintf("%s%s", userID, ext)
	filePath := filepath.Join(dir, filename)

	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return "/uploads/avatars/" + filename, nil
}
