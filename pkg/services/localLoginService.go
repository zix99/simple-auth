package services

import (
	"errors"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/lib/totp"
	"simple-auth/pkg/saerrors"
)

type LocalLoginService interface {
	FindAuthLocal(accountUUID string) (*db.AuthLocal, error)
	UsernameExists(username string) (bool, error)

	AssertLoginCredentialsOnly(username, password string) (*db.AuthLocal, error)
	AssertLogin(username, password string, totpCode *string) (*db.AuthLocal, error)

	ActivateTOTP(authLocal *db.AuthLocal, otp *totp.Totp, code string) error
	DeactivateTOTP(authLocal *db.AuthLocal, code string) error

	UpdatePassword(authLocal *db.AuthLocal, oldPassword string, newPassword string) error
	UpdatePasswordUnsafe(authLocal *db.AuthLocal, newPassword string) error
}

type localLoginService struct {
	dbAccount      db.AccountStore
	dbAuth         db.AccountAuthLocal
	dbAudit        db.AccountAudit
	dbStipulations db.AccountStipulations
	tfConfig       *config.TwoFactorConfig
}

var _ LocalLoginService = &localLoginService{}

func NewLocalLoginService(db db.SADB, tfConfig *config.TwoFactorConfig) LocalLoginService {
	return &localLoginService{
		dbAccount:      db,
		dbAuth:         db,
		dbAudit:        db,
		dbStipulations: db,
		tfConfig:       tfConfig,
	}
}

const (
	InvalidAccount               saerrors.ErrorCode = "invalid-account"
	LocalAuthMissing             saerrors.ErrorCode = "missing-local-auth"
	LocalInvalidCredentials      saerrors.ErrorCode = "invalid-credentials"
	LocalTOTPMissing             saerrors.ErrorCode = "totp-missing"
	LocalTOTPFailed              saerrors.ErrorCode = "totp-failed"
	LocalUnsatisfiedStipulations saerrors.ErrorCode = "unsatisfied-stipulations"
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
		return false, err
	}
	if localAuth != nil {
		return true, nil
	}
	return false, nil
}

func (s *localLoginService) AssertLogin(username, password string, totpCode *string) (*db.AuthLocal, error) {
	localAuth, err := s.dbAuth.FindAuthLocalByUsername(username)
	if err != nil {
		return nil, LocalInvalidCredentials.Wrap(err)
	}

	if err := s.assertLoginCredentials(localAuth, password); err != nil {
		return nil, err
	}

	if localAuth.HasTOTP() {
		if totpCode == nil || *totpCode == "" {
			return nil, LocalTOTPMissing.New()
		}
		if !localAuth.VerifyTOTP(*totpCode, s.tfConfig.Drift) {
			s.dbAudit.CreateAuditRecord(localAuth, db.AuditModuleLocal, db.AuditLevelWarn, "TOTP Rejected")
			return nil, LocalTOTPFailed.New()
		}
	}

	s.dbAudit.CreateAuditRecord(localAuth, db.AuditModuleLocal, db.AuditLevelInfo, "Login Successful")

	return localAuth, nil
}

func (s *localLoginService) AssertLoginCredentialsOnly(username, password string) (*db.AuthLocal, error) {
	localAuth, err := s.dbAuth.FindAuthLocalByUsername(username)
	if err != nil {
		return nil, LocalInvalidCredentials.Wrap(err)
	}

	if err := s.assertLoginCredentials(localAuth, password); err != nil {
		return nil, err
	}

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
	if !otp.Validate(verificationCode, s.tfConfig.Drift) {
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

	if !authLocal.VerifyTOTP(verificationCode, s.tfConfig.Drift) {
		return LocalTOTPFailed.New()
	}

	if err := s.dbAuth.UpdateAuthLocalTOTP(authLocal, nil); err != nil {
		return err
	}

	return nil
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
