package db

import (
	"github.com/labstack/gommon/random"
)

type TokenStipulation struct {
	Code string `json:"code"`
}

func (s *TokenStipulation) Type() StipulationType {
	return StipulationType("token")
}

func (s *TokenStipulation) IsSatisfiedBy(spec IStipulation) bool {
	other, ok := spec.(*TokenStipulation)
	if !ok {
		return false
	}
	return other.Code == s.Code
}

func NewTokenStipulation() *TokenStipulation {
	return &TokenStipulation{
		Code: random.String(18),
	}
}
