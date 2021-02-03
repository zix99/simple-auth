package services

import (
	"crypto/rand"
	"errors"
	"simple-auth/pkg/appcontext"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"time"

	"github.com/google/uuid"
)

type IssuedToken struct {
	AccessToken  string
	Expires      int
	RefreshToken string
	Scope        db.OAuthScope
}

type AuthOAuthService interface {
	WithContext(ctx appcontext.Context) AuthOAuthService
	CreateAccessCode(account *db.Account, scopes db.OAuthScope) (string, error)
	TradeCodeForToken(secret, code string) (ret IssuedToken, err error)
	TradeRefreshTokenForAccessToken(secret, refreshToken string) (ret IssuedToken, err error)
	TradeCredentialsForToken(secret, username, password string, totp *string, scopes db.OAuthScope) (ret IssuedToken, err error)

	FindExistingToken(account *db.Account, tokenType db.OAuthTokenType, scopes db.OAuthScope) (IssuedToken, error)

	ValidateRedirectURI(uri string) bool
	ValidateScopes(scopes db.OAuthScope) bool
}

var (
	ErrInvalidScopes = errors.New("invalid scope")
)

type authOAuthService struct {
	clientID string
	config   *config.ConfigOAuth2Client
	common   *config.ConfigOAuth2Settings

	dbOAuth    db.AccountOAuth
	localLogin LocalLoginService
}

func NewAuthOAuthService(clientID string, config *config.ConfigOAuth2Client, common *config.ConfigOAuth2Settings, localLoginService LocalLoginService) AuthOAuthService {
	return &authOAuthService{
		clientID,
		config,
		common,
		nil,
		localLoginService,
	}
}

func (s *authOAuthService) WithContext(ctx appcontext.Context) AuthOAuthService {
	copy := *s
	copy.dbOAuth = appcontext.GetSADB(ctx)
	copy.localLogin = s.localLogin.WithContext(ctx)
	return &copy
}

func (s *authOAuthService) CreateAccessCode(account *db.Account, scopes db.OAuthScope) (string, error) {
	if !s.ValidateScopes(scopes) {
		return "", ErrInvalidScopes
	}

	code, err := genAccessCode(s.common.CodeLength)
	if err != nil {
		return "", err
	}

	if err := s.dbOAuth.CreateOAuthToken(account, s.clientID, db.OAuthTypeCode, code, scopes, time.Duration(s.common.CodeExpiresSeconds)*time.Second); err != nil {
		return "", err
	}

	return code, nil
}

func (s *authOAuthService) TradeCodeForToken(secret, code string) (ret IssuedToken, err error) {
	if s.config.Secret != secret {
		err = errors.New("invalid secret")
		return
	}

	var token *db.OAuthToken
	token, err = s.dbOAuth.AssertOAuthToken(code, db.OAuthTypeCode, true)
	if err != nil {
		return
	}

	ret, err = s.issueToken(token.Account, token.Scopes)
	return
}

func (s *authOAuthService) TradeCredentialsForToken(secret, username, password string, totp *string, scopes db.OAuthScope) (ret IssuedToken, err error) {
	if !s.common.AllowCredentials {
		err = errors.New("trading credentials for token is disabled")
		return
	}
	if s.config.Secret != secret {
		err = errors.New("invalid secret")
		return
	}

	authLocal, err := s.localLogin.AssertLogin(username, password, totp)
	if err != nil {
		return
	}

	ret, err = s.issueToken(authLocal.Account(), scopes)
	return
}

func (s *authOAuthService) issueToken(account *db.Account, scopes db.OAuthScope) (ret IssuedToken, err error) {
	if !s.ValidateScopes(scopes) {
		err = ErrInvalidScopes
		return
	}

	if s.common.ReuseToken {
		if tokens, vtErr := s.dbOAuth.GetValidOAuthTokens(s.clientID, account); vtErr == nil {
			for _, t := range tokens {
				if ret.AccessToken == "" && t.Type == db.OAuthTypeAccessToken && t.Scopes.Matches(scopes) {
					ret.AccessToken = t.Token
					ret.Expires = int(time.Until(t.Expires).Seconds())
					ret.Scope = t.Scopes
				}
				if ret.RefreshToken == "" && t.Type == db.OAuthTypeRefreshToken && s.config.IssueRefreshToken {
					ret.RefreshToken = t.Token
				}
			}
			if ret.AccessToken != "" && (ret.RefreshToken != "" || !s.config.IssueRefreshToken) {
				return
			}
		}
	}

	s.dbOAuth.InvalidateAllOAuth(s.clientID, account)

	ret.AccessToken = uuid.New().String()
	ret.Expires = s.common.TokenExpiresSeconds
	ret.Scope = scopes
	err = s.dbOAuth.CreateOAuthToken(account, s.clientID, db.OAuthTypeAccessToken, ret.AccessToken, scopes, time.Duration(s.common.TokenExpiresSeconds)*time.Second)
	if err != nil {
		return
	}

	if s.config.IssueRefreshToken {
		ret.RefreshToken = uuid.New().String()
		const oneHundredYears = 100 * 365 * 24 * time.Hour
		err = s.dbOAuth.CreateOAuthToken(account, s.clientID, db.OAuthTypeRefreshToken, ret.RefreshToken, scopes, oneHundredYears)
		if err != nil {
			return
		}
	}

	return
}

func (s *authOAuthService) TradeRefreshTokenForAccessToken(secret, refreshToken string) (ret IssuedToken, err error) {
	if s.config.Secret != secret {
		err = errors.New("invalid secret")
		return
	}

	var token *db.OAuthToken
	token, err = s.dbOAuth.AssertOAuthToken(refreshToken, db.OAuthTypeRefreshToken, false)
	if err != nil {
		return
	}

	ret.AccessToken = uuid.New().String()
	ret.Expires = s.common.TokenExpiresSeconds
	ret.Scope = token.Scopes
	err = s.dbOAuth.CreateOAuthToken(token.Account, s.clientID, db.OAuthTypeAccessToken, ret.AccessToken, token.Scopes, time.Duration(s.common.TokenExpiresSeconds)*time.Second)
	if err != nil {
		return
	}

	return
}

func (s *authOAuthService) FindExistingToken(account *db.Account, tokenType db.OAuthTokenType, scopes db.OAuthScope) (IssuedToken, error) {
	tokens, err := s.dbOAuth.GetValidOAuthTokens(s.clientID, account)
	if err != nil {
		return IssuedToken{}, err
	}

	for _, v := range tokens {
		if v.Type == tokenType && v.Scopes.Matches(scopes) {
			return IssuedToken{
				AccessToken: v.Token,
				Expires:     int(time.Until(v.Expires).Seconds()),
				Scope:       v.Scopes,
			}, nil
		}
	}

	return IssuedToken{}, errors.New("no token found")
}

func (s *authOAuthService) ValidateRedirectURI(uri string) bool {
	return uri == s.config.RedirectURI
}

func (s *authOAuthService) ValidateScopes(scopes db.OAuthScope) bool {
	return db.OAuthScope(s.config.Scopes).ContainsAll(scopes...)
}

func (s *authOAuthService) InspectToken(sToken string) (*db.OAuthToken, error) {
	token, err := s.dbOAuth.AssertOAuthToken(sToken, db.OAuthTypeAccessToken, false)
	return token, err
}

func genAccessCode(digits int) (string, error) {
	var table = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	ret := make([]byte, digits)
	if _, err := rand.Read(ret); err != nil {
		return "", err
	}

	for i := 0; i < digits; i++ {
		ret[i] = table[int(ret[i])%len(table)]
	}
	return string(ret), nil
}
