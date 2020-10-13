package db

import "simple-auth/pkg/saerrors"

// Common
const (
	// Common
	InternalError   saerrors.ErrorCode = "db-internal-error"
	InvalidAccount  saerrors.ErrorCode = "invalid-account"
	InactiveAccount saerrors.ErrorCode = "inactive-account"

	// authOneTime
	SAOneTimeInvalidToken saerrors.ErrorCode = "invalid-token"
	SAOneTimeConsumed     saerrors.ErrorCode = "consumed"
	SAOneTimeExpired      saerrors.ErrorCode = "expired"

	// authSimple
	SAInvalidCredentials      saerrors.ErrorCode = "invalid-credentials"
	SAUserVerificationFailed  saerrors.ErrorCode = "user-verification-failed"
	SAUnsatisfiedStipulations saerrors.ErrorCode = "unsatisfied-stipulations"
	SATOTPMissing             saerrors.ErrorCode = "totp-missing"
	SATOTPFailed              saerrors.ErrorCode = "totp-failed"

	// authToken
	SessionNotFound      saerrors.ErrorCode = "session-missing"
	SessionInvalidated   saerrors.ErrorCode = "session-invalid"
	SessionExpired       saerrors.ErrorCode = "session-expired"
	VerificationMissing  saerrors.ErrorCode = "verification-missing"
	VerificationConsumed saerrors.ErrorCode = "verification-consumed"
	VerificationExpired  saerrors.ErrorCode = "verification-expired"
	VerificationInvalid  saerrors.ErrorCode = "verification-invalid"
)
