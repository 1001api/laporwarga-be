package auth

import (
	"errors"
	"fmt"
	db "hubku/lapor_warga_be_v2/internal/database/generated"
	"hubku/lapor_warga_be_v2/internal/modules/users"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(req LoginRequest) (*LoginResponse, error)
	ValidateToken(tokenString string) (*Claims, error)
	GenerateToken(user db.GetUserByIdentifierRow) (string, error)
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error
	Register(req RegisterRequest) (db.CreateUserRow, error)
	GenerateRefreshToken(user db.GetUserByIdentifierRow) (string, error)
	RefreshToken(req RefreshRequest) (*LoginResponse, error)
}

type service struct {
	userService   users.UserService
	jwtSecret     []byte
	tokenExpiry   time.Duration
	refreshExpiry time.Duration
}

func NewAuthService(userService users.UserService) AuthService {
	viper.SetDefault("JWT_EXPIRY", 15)
	viper.SetDefault("JWT_REFRESH_EXPIRY", 720)

	return &service{
		userService:   userService,
		jwtSecret:     []byte(viper.GetString("JWT_SECRET")),
		tokenExpiry:   time.Duration(viper.GetInt("JWT_EXPIRY")) * time.Minute,
		refreshExpiry: time.Duration(viper.GetInt("JWT_REFRESH_EXPIRY")) * time.Minute,
	}
}

func (s *service) Login(req LoginRequest) (*LoginResponse, error) {
	if req.Identifier == "" || req.Password == "" {
		return nil, errors.New("username and password are required")
	}

	user, err := s.userService.GetUserByIdentifier(req.Identifier)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if cast.ToTime(user.LockedUntil).After(time.Now()) {
		return nil, errors.New("account is temporarily locked due to multiple failed login attempts, please wait a few minutes")
	}

	if err := s.VerifyPassword(user.PasswordHash.String, req.Password); err != nil {
		s.userService.IncrementFailedLogins(user.ID)
		log.Println(err)
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := s.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}
	refreshToken, err := s.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// update last login in background
	go func(userID uuid.UUID) {
		s.userService.UpdateUserLastLogin(userID)
	}(user.ID)

	return &LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) Register(req RegisterRequest) (db.CreateUserRow, error) {
	hashed, err := s.HashPassword(req.Password)
	if err != nil {
		return db.CreateUserRow{}, err
	}

	params := db.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		FullName:     req.Fullname,
		PasswordHash: hashed,
		PhoneNumber:  req.PhoneNumber,
		Role:         req.Role,
	}

	return s.userService.CreateUser(params)
}

func (s *service) GenerateToken(user db.GetUserByIdentifierRow) (string, error) {
	expiresAt := time.Now().Add(s.tokenExpiry)

	claims := &Claims{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role.String,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "lapor_warga",
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *service) GenerateRefreshToken(user db.GetUserByIdentifierRow) (string, error) {
	expiresAt := time.Now().Add(s.refreshExpiry)

	claims := &Claims{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role.String,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ims_be_v1",
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *service) RefreshToken(req RefreshRequest) (*LoginResponse, error) {
	claims, err := s.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid token type")
	}

	user, err := s.userService.GetUserByIdentifier(claims.UserID.String())
	if err != nil {
		return nil, errors.New("user not found")
	}

	accessToken, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:        accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *service) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hash), err
}

func (s *service) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return s.jwtSecret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
