package auth

import (
	"errors"
	"fmt"
	db "hubku/lapor_warga_be_v2/internal/database/generated"
	"hubku/lapor_warga_be_v2/internal/modules/auditlogs"
	"hubku/lapor_warga_be_v2/internal/modules/users"
	"hubku/lapor_warga_be_v2/pkg"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type AuthService interface {
	Login(req LoginRequest, isMobile bool) (*LoginResponse, error)
	ValidateToken(tokenString string) (*Claims, error)
	GenerateToken(user db.GetUserByIdentifierRow) (string, error)
	Register(currentUserID uuid.UUID, req RegisterRequest) (uuid.UUID, error)
	GenerateRefreshToken(user db.GetUserByIdentifierRow) (string, error)
	RefreshToken(req RefreshRequest) (*LoginResponse, error)
}

type service struct {
	userService   users.UserService
	logService    auditlogs.LogsService
	jwtSecret     []byte
	tokenExpiry   time.Duration
	refreshExpiry time.Duration
	enckey        []byte
}

func NewAuthService(userService users.UserService, logService auditlogs.LogsService, encKey string) AuthService {
	viper.SetDefault("JWT_EXPIRY", 15)
	viper.SetDefault("JWT_REFRESH_EXPIRY", 720)

	return &service{
		userService:   userService,
		logService:    logService,
		jwtSecret:     []byte(viper.GetString("JWT_SECRET")),
		tokenExpiry:   time.Duration(viper.GetInt("JWT_EXPIRY")) * time.Minute,
		refreshExpiry: time.Duration(viper.GetInt("JWT_REFRESH_EXPIRY")) * time.Minute,
		enckey:        []byte(encKey),
	}
}

func (s *service) Login(req LoginRequest, isMobile bool) (*LoginResponse, error) {
	if req.Identifier == "" || req.Password == "" {
		return nil, errors.New("username and password are required")
	}

	user, err := s.userService.GetUserByIdentifier(req.Identifier)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// check if user role is citizen
	// this is to prevent citizen from logging in using web route.
	if isMobile {
		if user.RoleName.String != string(pkg.RoleCitizen) {
			return nil, errors.New("only citizen can login to this route")
		}
	} else {
		if user.RoleName.String == string(pkg.RoleCitizen) {
			return nil, errors.New("citizen cannot login to this route")
		}
	}

	if cast.ToTime(user.LockedUntil).After(time.Now()) {
		return nil, errors.New("account is temporarily locked due to multiple failed login attempts, please wait a few minutes")
	}

	if err := pkg.VerifyPassword(user.PasswordHash.String, req.Password); err != nil {
		s.userService.IncrementFailedLogins(user.ID)
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

	go func(userID uuid.UUID) {
		// update last login in background
		s.userService.UpdateUserLastLogin(userID)

		// create log in background
		s.logService.CreateLog(db.CreateAuditLogParams{
			EntityName:  string(pkg.LogEntityUsers),
			Action:      string(pkg.LogTypeLogin),
			EntityID:    user.ID,
			PerformedBy: user.ID,
		})
	}(user.ID)

	return &LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) Register(currentUserID uuid.UUID, req RegisterRequest) (uuid.UUID, error) {
	hashed, err := pkg.HashPassword(req.Password)
	if err != nil {
		return uuid.UUID{}, err
	}

	createdID, err := s.userService.CreateUser(currentUserID, users.CreateUserRequest{
		Username:     req.Username,
		Email:        req.Email,
		FullName:     req.Fullname,
		PasswordHash: string(hashed),
		PhoneNumber:  req.PhoneNumber,
		Role:         req.Role,
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	return createdID, nil
}

func (s *service) GenerateToken(user db.GetUserByIdentifierRow) (string, error) {
	expiresAt := time.Now().Add(s.tokenExpiry)

	claims := &Claims{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     string(user.Email),
		Role:      user.RoleName.String,
		TokenType: pkg.AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    pkg.JWTIssuer,
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
		Email:     string(user.Email),
		Role:      user.RoleName.String,
		TokenType: pkg.RefreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    pkg.JWTIssuer,
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

	if claims.TokenType != pkg.RefreshToken {
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
