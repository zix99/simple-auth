package services

import "simple-auth/pkg/config"

type OIDCService interface {
	GetProvider(id string) *config.ConfigOIDCProvider
}

type oidcService struct {
	providers []*config.ConfigOIDCProvider
}

func NewOIDCService(providers []*config.ConfigOIDCProvider) OIDCService {
	return &oidcService{
		providers,
	}
}

func (s *oidcService) GetProvider(id string) *config.ConfigOIDCProvider {
	if id == "" {
		return nil
	}
	if s.providers == nil {
		return nil
	}
	for _, provider := range s.providers {
		if provider.ID == id {
			return provider
		}
	}
	return nil
}
