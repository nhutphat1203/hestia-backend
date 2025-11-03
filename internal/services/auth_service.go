package service

import (
	"errors"

	"github.com/nhutphat1203/hestia-backend/internal/config"
	"github.com/nhutphat1203/hestia-backend/internal/infrastructure/auth"
	repository "github.com/nhutphat1203/hestia-backend/internal/repositories"
	"github.com/nhutphat1203/hestia-backend/pkg/gen"
	hasher "github.com/nhutphat1203/hestia-backend/pkg/hash"
	"github.com/nhutphat1203/hestia-backend/pkg/logger"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthService struct {
	userRepo        repository.UserRepository
	userSessionRepo repository.UserSessionRepository
	cfg             *config.Config
	jwtService      *auth.JWTService
	logger          *logger.Logger
}

func NewAuthService(userRepo repository.UserRepository, userSessionRepo repository.UserSessionRepository, cfg *config.Config, logger *logger.Logger) *AuthService {
	jwtService := auth.NewJWTService(cfg.JWTSecret, "access_token", int(cfg.JWTExpiration.Minutes()))
	return &AuthService{
		userRepo:        userRepo,
		userSessionRepo: userSessionRepo,
		cfg:             cfg,
		jwtService:      jwtService,
		logger:          logger,
	}
}

func (s *AuthService) Login(account, password string) (TokenResponse, error) {
	user, err := s.userRepo.GetUserByAccount(account)
	if err != nil {
		return TokenResponse{}, err
	}
	if !hasher.Verify(password, user.HashedPassword) {
		return TokenResponse{}, errors.New("invalid credentials")
	}

	accessToken, err := s.jwtService.GenerateToken(int(user.ID))
	if err != nil {
		return TokenResponse{}, err
	}

	refreshToken := gen.GenerateToken()

	hashedToken, err := hasher.Hash(refreshToken)
	if err != nil {
		return TokenResponse{}, err
	}

	latestSession, err := s.userSessionRepo.GetLatestSessionByUserId(user.ID)

	if (err == nil) && (!latestSession.IsRevoked) {
		latestSession.IsRevoked = true
		_ = s.userSessionRepo.Update(latestSession)
	}

	_, err = s.userSessionRepo.Create(user.ID, hashedToken)

	if err != nil {
		return TokenResponse{}, err
	}

	return TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	hashedToken, err := hasher.Hash(refreshToken)
	if err != nil {
		return err
	}
	userSession, err := s.userSessionRepo.GetSessionByToken(hashedToken)
	if err != nil {
		return err
	}
	userSession.IsRevoked = true
	return s.userSessionRepo.Update(userSession)
}

func (s *AuthService) RefreshToken(userId uint, refreshToken string) (TokenResponse, error) {
	hashedToken, err := hasher.Hash(refreshToken)
	if err != nil {
		return TokenResponse{}, err
	}
	userSession, err := s.userSessionRepo.GetSessionByUserIdAndToken(userId, hashedToken)

	if err != nil {
		return TokenResponse{}, err
	}

	if userSession.IsRevoked {
		return TokenResponse{}, errors.New("invalid token")
	}

	accessToken, err := s.jwtService.GenerateToken(int(userSession.UserID))

	if err != nil {
		return TokenResponse{}, err
	}
	newRefreshToken := gen.GenerateToken()

	newHashedToken, err := hasher.Hash(newRefreshToken)

	if err != nil {
		return TokenResponse{}, err
	}

	s.Logout(refreshToken)

	_, err = s.userSessionRepo.Create(userSession.UserID, newHashedToken)

	if err != nil {
		return TokenResponse{}, err
	}

	return TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
