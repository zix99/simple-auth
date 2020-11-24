package v1

import (
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/services"
)

type Environment struct {
	localLoginService services.LocalLoginService
}

func NewEnvironment(config *config.ConfigWeb, db db.SADB) *Environment {
	return &Environment{
		services.NewLocalLoginService(db, &config.Login.TwoFactor),
	}
}
