package auth_usecase

import (
	"core/config"
	"core/internal/auth"
	"core/internal/failure"
	"core/pkg/secure"
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

	correctSignature := secure.CalcSignature(secret.ApiPrivate, secret.ApiPublic+timestamp+requestId+body)

	return signature == correctSignature, nil
}
