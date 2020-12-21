package services

import (
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/lib/totp"
)

type TwoFactorService interface {
	CreateSecret() (string, error)
	CreateFullSpecFromSecret(secret string, authLocal *db.AuthLocal) (*totp.Totp, error)
}

type twoFactorService struct {
	tfConfig *config.ConfigTwoFactor
}

var _ TwoFactorService = &twoFactorService{}

func NewTwoFactorService(tfConfig *config.ConfigTwoFactor) TwoFactorService {
	return &twoFactorService{
		tfConfig,
	}
}

func (s *twoFactorService) CreateSecret() (string, error) {
	secret, err := totp.CreateSecret(s.tfConfig.KeyLength)
	if err != nil {
		return "", err
	}

	return totp.EncodeSecretb32(secret), nil
}

func (s *twoFactorService) CreateFullSpecFromSecret(secret string, authLocal *db.AuthLocal) (*totp.Totp, error) {
	t, err := totp.FromSecret(secret, s.tfConfig.Issuer, authLocal.Account().Email)
	if err != nil {
		return nil, err
	}
	return t, nil
}
