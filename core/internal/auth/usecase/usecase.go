package auth_usecase

import (
	"market_auth/config"
	"market_auth/internal/auth"
	"market_auth/internal/failure"
	"market_auth/pkg/secure"
)

type uc struct {
	cfg *config.Config
}

func New(cfg *config.Config) auth.UC {
	return &uc{
		cfg: cfg,
	}
}

func (u *uc) VerifySignature(service, signature, body, timestamp, requestId string) (bool, error) {
	secret, ok := u.cfg.Secrets[service]
	if !ok {
		return false, failure.ErrAuth
	}

	return signature == secure.CalcSignature(secret.ApiPrivate, secret.ApiPublic+timestamp+requestId+body), nil
}
