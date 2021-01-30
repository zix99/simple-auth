package db

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type AccountOAuth interface {
	CreateOAuthToken(account *Account, clientID string, tokenType OAuthTokenType, token string, scopes OAuthScope, expiresIn time.Duration) error
	AssertOAuthToken(token string, tokenType OAuthTokenType, consume bool) (*OAuthToken, error)
	InvalidateAllOAuth(clientId string, account *Account) error

	GetValidOAuthTokens(clientID string, account *Account) ([]*OAuthToken, error)
	GetAllValidOAuthTokens(account *Account) ([]*OAuthToken, error)
}

type OAuthTokenType string

const (
	OAuthTypeRefreshToken OAuthTokenType = "refresh_token"
	OAuthTypeAccessToken  OAuthTokenType = "access_token"
	OAuthTypeCode         OAuthTokenType = "code"
)

type accountOAuthToken struct {
	gorm.Model
	AccountID   uint `gorm:"index; not null"`
	ClientID    string
	Type        OAuthTokenType
	Token       string `gorm:"uniqueIndex; not null"`
	Scope       string
	Expires     time.Time
	Invalidated bool
}

type OAuthToken struct {
	Account  *Account
	Scopes   OAuthScope
	Token    string
	ClientID string
	Type     OAuthTokenType
	Created  time.Time
	Expires  time.Time
}

func (s *sadb) CreateOAuthToken(account *Account, clientID string, tokenType OAuthTokenType, token string, scopes OAuthScope, expiresIn time.Duration) error {
	if account == nil || clientID == "" || tokenType == "" || token == "" {
		return errors.New("invalid params")
	}

	oauth := &accountOAuthToken{
		AccountID:   account.ID,
		ClientID:    clientID,
		Type:        tokenType,
		Token:       token,
		Scope:       scopes.String(),
		Expires:     time.Now().Add(expiresIn),
		Invalidated: false,
	}

	if err := s.db.Create(oauth).Error; err != nil {
		return err
	}

	s.CreateAuditRecord(account, AuditModuleOAuth2, AuditLevelInfo, "Issued OAuth2 %s", tokenType)

	return nil
}

func (s *sadb) AssertOAuthToken(token string, tokenType OAuthTokenType, consume bool) (*OAuthToken, error) {
	if token == "" {
		return nil, errors.New("no token")
	}
	if tokenType == "" {
		return nil, errors.New("no token type")
	}

	var oauth accountOAuthToken
	if err := s.db.Where("token = ? AND type = ?", token, tokenType).First(&oauth).Error; err != nil {
		return nil, err
	}

	if oauth.Invalidated {
		return nil, errors.New("invalidated token")
	}
	if time.Now().After(oauth.Expires) {
		return nil, errors.New("expired token")
	}

	var account Account
	if err := s.db.Model(&oauth).Related(&account).Error; err != nil {
		return nil, err
	}

	if !account.Active {
		return nil, errors.New("inactive account")
	}

	if consume {
		err := s.db.Model(&oauth).Update(&accountOAuthToken{
			Invalidated: true,
		}).Error
		if err != nil {
			return nil, err
		}
	}

	return dbTokenToOAuthToken(&account, &oauth), nil
}

func (s *sadb) GetValidOAuthTokens(clientID string, account *Account) ([]*OAuthToken, error) {
	if clientID == "" {
		return nil, errors.New("invalid client")
	}
	if account == nil {
		return nil, errors.New("invalid account")
	}

	var tokens []*accountOAuthToken
	if err := s.db.Where("account_id = ? AND client_id = ? AND invalidated = false", account.ID, clientID).Find(&tokens).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	ret := make([]*OAuthToken, 0, len(tokens))
	for _, token := range tokens {
		if !token.Invalidated && token.Expires.After(now) {
			ret = append(ret, dbTokenToOAuthToken(account, token))
		}
	}
	return ret, nil
}

func (s *sadb) GetAllValidOAuthTokens(account *Account) ([]*OAuthToken, error) {
	if account == nil {
		return nil, errors.New("invalid account")
	}

	var tokens []*accountOAuthToken
	if err := s.db.Where("account_id = ? AND invalidated = false", account.ID).Find(&tokens).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	ret := make([]*OAuthToken, 0, len(tokens))
	for _, token := range tokens {
		if !token.Invalidated && token.Expires.After(now) {
			ret = append(ret, dbTokenToOAuthToken(account, token))
		}
	}
	return ret, nil
}

func (s *sadb) InvalidateAllOAuth(clientId string, account *Account) error {
	if clientId == "" {
		return errors.New("invalid clientId")
	}
	if account == nil {
		return errors.New("invalid account")
	}
	return s.db.Model(&accountOAuthToken{}).Where("client_id = ? and account_id = ?", clientId, account.ID).Update("invalidated", true).Error
}

func dbTokenToOAuthToken(account *Account, token *accountOAuthToken) *OAuthToken {
	return &OAuthToken{
		account,
		NewOAuthScope(token.Scope),
		token.Token,
		token.ClientID,
		token.Type,
		token.CreatedAt,
		token.Expires,
	}
}
