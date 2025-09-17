package auth

import (
	"errors"

	"github.com/nhutphat1203/hestia-backend/internal/domain"
)

var _ domain.Authenticator = (*StaticTokenAuth)(nil)

type StaticTokenAuth struct {
	Token string
}

func (s *StaticTokenAuth) Authenticate(token string) (bool, error) {
	if token == "" {
		return false, errors.New("missing token")
	}
	if token != s.Token {
		return false, errors.New("invalid token")
	}
	return true, nil
}

func NewStaticTokenAuth(token string) *StaticTokenAuth {
	return &StaticTokenAuth{
		Token: token,
	}
}
