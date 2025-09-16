package domain

type Authenticator interface {
	Authenticate(token string) (bool, error)
}
