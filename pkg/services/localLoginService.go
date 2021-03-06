package services

import (
	"errors"
	"fmt"
	"html/template"
	"regexp"
	"simple-auth/pkg/appcontext"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/email"
	"simple-auth/pkg/lib/totp"
	"simple-auth/pkg/saerrors"
	"unicode/utf8"
)

type LocalLoginService interface {
	FindAuthLocal(accountUUID string) (*db.AuthLocal, error)
	Create(account *db.Account, username, password string) (*db.AuthLocal, error)
	UsernameExists(username string) (bool, error)

	AssertLogin(usernameOrEmail, password string, totpCode *string) (*db.AuthLocal, error)

	ActivateTOTP(authLocal *db.AuthLocal, otp *totp.Totp, code string) error
	DeactivateTOTP(authLocal *db.AuthLocal, code string) error
	AllowTOTP() bool

	UpdatePassword(authLocal *db.AuthLocal, oldPassword string, newPassword string) error
	UpdatePasswordUnsafe(authLocal *db.AuthLocal, newPassword string) error

	WithContext(ctx appcontext.Context) LocalLoginService
}
type localLoginService struct {
	dbAccount      db.AccountStore
	dbAuth         db.AccountAuthLocal
	dbAudit        db.AccountAudit
	dbStipulations db.AccountStipulations
	emailService   *email.EmailService
	metaConfig     *config.ConfigMetadata
	lpConfig       *config.ConfigLocalProvider
	baseURL        string
}

var _ LocalLoginService = &localLoginService{}

func NewLocalLoginService(emailService *email.EmailService, metaConfig *config.ConfigMetadata, localProviderConfig *config.ConfigLocalProvider, baseURL string) LocalLoginService {
	return &localLoginService{
		emailService: emailService,
		metaConfig:   metaConfig,
		lpConfig:     localProviderConfig,
		baseURL:      baseURL,
	}
}

func (s *localLoginService) WithContext(ctx appcontext.Context) LocalLoginService {
	copy := *s
	db := appcontext.GetSADB(ctx)
	copy.dbAccount = db
	copy.dbAudit = db
	copy.dbAuth = db
	copy.dbStipulations = db
	return &copy
}

const (
	InvalidAccount               saerrors.ErrorCode = "invalid-account"
	LocalAuthMissing             saerrors.ErrorCode = "missing-local-auth"
	LocalInvalidCredentials      saerrors.ErrorCode = "invalid-credentials"
	LocalTOTPMissing             saerrors.ErrorCode = "totp-missing"
	LocalTOTPFailed              saerrors.ErrorCode = "totp-failed"
	LocalUnsatisfiedStipulations saerrors.ErrorCode = "unsatisfied-stipulations"
	LocalCredentialRequirements  saerrors.ErrorCode = "credentials-failed-requirements"
	LocalUsernameUnavailable     saerrors.ErrorCode = "username-unavailable"
)

func (s *localLoginService) FindAuthLocal(accountUUID string) (*db.AuthLocal, error) {
	account, err := s.dbAccount.FindAccount(accountUUID)
	if err != nil {
		return nil, InvalidAccount.Wrap(err)
	}

	authLocal, err := s.dbAuth.FindAuthLocal(account)
	if err != nil {
		return nil, LocalAuthMissing.Wrap(err)
	}

	return authLocal, nil
}

func (s *localLoginService) UsernameExists(username string) (bool, error) {
	localAuth, err := s.dbAuth.FindAuthLocalByUsername(username)
	if err != nil {
		return false, nil // FIXME: This is always an error-case the way the db func works
	}
	if localAuth != nil {
		return true, nil
	}
	return false, nil
}

func (s *localLoginService) Create(account *db.Account, username, password string) (*db.AuthLocal, error) {
	if err := s.validateUsername(username); err != nil {
		return nil, LocalCredentialRequirements.Compose(err)
	}
	if err := s.validatePassword(password); err != nil {
		return nil, LocalCredentialRequirements.Compose(err)
	}

	if authLocal, _ := s.dbAuth.FindAuthLocalByUsername(username); authLocal != nil {
		return nil, LocalUsernameUnavailable.New()
	}

	authLocal, err := s.dbAuth.CreateAuthLocal(account, username, password)
	if err != nil {
		return nil, err
	}

	if s.lpConfig.EmailValidationRequired {
		stip := db.NewTokenStipulation()
		s.dbStipulations.AddStipulation(account, stip)

		go s.emailService.SendVerificationEmail(account.Email, &email.VerificationData{
			EmailData: email.EmailData{
				Company: s.metaConfig.Company,
				BaseURL: s.baseURL,
			},
			ActivationLink: template.HTML(fmt.Sprintf("%s/#/activate?account=%s&token=%s", s.baseURL, account.UUID, stip.Code)),
		})
	}

	return authLocal, err
}

func (s *localLoginService) validateUsername(username string) error {
	ulen := utf8.RuneCountInString(username)
	if ulen < s.lpConfig.Requirements.UsernameMinLength {
		return errors.New("username too short")
	}
	if ulen > s.lpConfig.Requirements.UsernameMaxLength {
		return errors.New("username too long")
	}

	if s.lpConfig.Requirements.UsernameRegex != "" {
		re, err := regexp.Compile(s.lpConfig.Requirements.UsernameRegex)
		if err != nil {
			return errors.New("unable to parse valid username regex, ask your server admin to fix this")
		}
		if !re.MatchString(username) {
			return errors.New("invalid username characters")
		}
	}

	return nil
}

func (s *localLoginService) validatePassword(password string) error {
	plen := utf8.RuneCountInString(password)
	if plen < s.lpConfig.Requirements.PasswordMinLength {
		return errors.New("password too short")
	}
	if plen > s.lpConfig.Requirements.PasswordMaxLength {
		return errors.New("password too long")
	}
	return nil
}

func (s *localLoginService) AssertLogin(usernameOrEmail, password string, totpCode *string) (*db.AuthLocal, error) {
	localAuth, err := s.dbAuth.FindAuthLocalByEmail(usernameOrEmail)
	if err != nil {
		localAuth, err = s.dbAuth.FindAuthLocalByUsername(usernameOrEmail)
		if err != nil {
			return nil, LocalInvalidCredentials.Wrap(err)
		}
	}

	if err := s.assertLoginCredentials(localAuth, password); err != nil {
		return nil, err
	}

	if localAuth.HasTOTP() {
		if totpCode == nil || *totpCode == "" {
			return nil, LocalTOTPMissing.New()
		}
		if !localAuth.VerifyTOTP(*totpCode, s.lpConfig.TwoFactor.Drift) {
			s.dbAudit.CreateAuditRecord(localAuth, db.AuditModuleLocal, db.AuditLevelWarn, "TOTP Rejected")
			return nil, LocalTOTPFailed.New()
		}
	}

	s.dbAudit.CreateAuditRecord(localAuth, db.AuditModuleLocal, db.AuditLevelInfo, "Login Successful")

	return localAuth, nil
}

func (s *localLoginService) assertLoginCredentials(localAuth *db.AuthLocal, password string) error {
	if !localAuth.Account().Active {
		return db.InactiveAccount.New()
	}

	if !localAuth.VerifyPassword(password) {
		s.dbAudit.CreateAuditRecord(localAuth, db.AuditModuleLocal, db.AuditLevelWarn, "Login failed")
		return LocalInvalidCredentials.New()
	}

	if s.dbStipulations.AccountHasUnsatisfiedStipulations(localAuth.Account()) {
		return LocalUnsatisfiedStipulations.New()
	}

	return nil
}

func (s *localLoginService) ActivateTOTP(authLocal *db.AuthLocal, otp *totp.Totp, verificationCode string) error {
	if !otp.Validate(verificationCode, s.lpConfig.TwoFactor.Drift) {
		return LocalTOTPFailed.New()
	}

	tStr := otp.String()
	if err := s.dbAuth.UpdateAuthLocalTOTP(authLocal, &tStr); err != nil {
		return err
	}

	return nil
}

func (s *localLoginService) DeactivateTOTP(authLocal *db.AuthLocal, verificationCode string) error {
	if !authLocal.HasTOTP() {
		return errors.New("totp disabled")
	}

	if !authLocal.VerifyTOTP(verificationCode, s.lpConfig.TwoFactor.Drift) {
		return LocalTOTPFailed.New()
	}

	if err := s.dbAuth.UpdateAuthLocalTOTP(authLocal, nil); err != nil {
		return err
	}

	return nil
}

func (s *localLoginService) AllowTOTP() bool {
	return s.lpConfig.TwoFactor.Enabled
}

func (s *localLoginService) UpdatePassword(authLocal *db.AuthLocal, oldPassword string, newPassword string) error {
	if !authLocal.VerifyPassword(oldPassword) {
		return LocalInvalidCredentials.New()
	}

	return s.UpdatePasswordUnsafe(authLocal, newPassword)
}

func (s *localLoginService) UpdatePasswordUnsafe(authLocal *db.AuthLocal, newPassword string) error {
	return s.dbAuth.UpdateAuthLocalPassword(authLocal, newPassword)
}
