package db

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type AccountOAuth interface {
	CreateOAuthToken(account *Account, clientID string, tokenType OAuthTokenType, token string, scopes OAuthScope, expiresIn time.Duration) error
	AssertOAuthToken(clientID, token string, tokenType OAuthTokenType, consume bool) (*OAuthToken, error)
	InvalidateToken(clientId string, account *Account, token string) error
	InvalidateAllOAuth(clientId string, account *Account, exceptType []OAuthTokenType) error

	// Missing will return nil,nil
	GetValidOAuthToken(token string) (*OAuthToken, error)

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
	AccountID uint `gorm:"index; not null"`
	ClientID  string
	Type      OAuthTokenType
	Token     string `gorm:"uniqueIndex; not null"`
	Scope     string
	Expires   time.Time
}

func (s *accountOAuthToken) Expired() bool {
	return time.Now().After(s.Expires)
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

func (s *OAuthToken) Expired() bool {
	return time.Now().After(s.Expires)
}

func (s *sadb) CreateOAuthToken(account *Account, clientID string, tokenType OAuthTokenType, token string, scopes OAuthScope, expiresIn time.Duration) error {
	if account == nil || clientID == "" || tokenType == "" || token == "" {
		return errors.New("invalid params")
	}

	oauth := &accountOAuthToken{
		AccountID: account.ID,
		ClientID:  clientID,
		Type:      tokenType,
		Token:     token,
		Scope:     scopes.String(),
		Expires:   time.Now().Add(expiresIn),
	}

	if err := s.db.Create(oauth).Error; err != nil {
		return err
	}

	s.CreateAuditRecord(account, AuditModuleOAuth2, AuditLevelInfo, "Issued OAuth2 %s", tokenType)

	return nil
}

func (s *sadb) AssertOAuthToken(clientID, token string, tokenType OAuthTokenType, consume bool) (*OAuthToken, error) {
	if token == "" {
		return nil, errors.New("no token")
	}
	if tokenType == "" {
		return nil, errors.New("no token type")
	}
	if clientID == "" {
		return nil, errors.New("no client_id")
	}

	var oauth accountOAuthToken
	if err := s.db.Where("token = ? AND type = ? AND client_id = ?", token, tokenType, clientID).First(&oauth).Error; err != nil {
		return nil, err
	}

	if oauth.Expired() {
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
		if err := s.db.Delete(&oauth).Error; err != nil {
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
	if err := s.db.Where("account_id = ? AND client_id = ?", account.ID, clientID).Find(&tokens).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	ret := make([]*OAuthToken, 0, len(tokens))
	for _, token := range tokens {
		if token.Expires.After(now) {
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
	if err := s.db.Where("account_id = ?", account.ID).Find(&tokens).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	ret := make([]*OAuthToken, 0, len(tokens))
	for _, token := range tokens {
		if token.Expires.After(now) {
			ret = append(ret, dbTokenToOAuthToken(account, token))
		}
	}
	return ret, nil
}

func (s *sadb) GetValidOAuthToken(token string) (*OAuthToken, error) {
	if token == "" {
		return nil, errors.New("missing token")
	}

	var oauth accountOAuthToken
	if err := s.db.Where("token = ?", token).Find(&oauth).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if oauth.Expired() {
		return nil, nil
	}

	var account Account
	if err := s.db.Model(&oauth).Related(&account).Error; err != nil {
		return nil, err
	}

	return dbTokenToOAuthToken(&account, &oauth), nil
}

func (s *sadb) InvalidateToken(clientId string, account *Account, token string) error {
	if clientId == "" {
		return errors.New("invalid clientId")
	}
	if account == nil {
		return errors.New("invalid account")
	}
	if token == "" {
		return errors.New("invalid token")
	}

	return s.db.Where("client_id = ? and account_id = ? and token = ?", clientId, account.ID, token).Delete(&accountOAuthToken{}).Error
}

func (s *sadb) InvalidateAllOAuth(clientId string, account *Account, exceptType []OAuthTokenType) error {
	if clientId == "" {
		return errors.New("invalid clientId")
	}
	if account == nil {
		return errors.New("invalid account")
	}

	q := s.db.Where("client_id = ? and account_id = ?", clientId, account.ID)
	if len(exceptType) > 0 {
		q = q.Not("type in (?)", exceptType)
	}

	return q.Delete(&accountOAuthToken{}).Error
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
