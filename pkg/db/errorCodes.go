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
	SAOneTimeExpired      saerrors.ErrorCode = "expired"

	// authLocal
	AuthInvalidUsername saerrors.ErrorCode = "invalid-username"

	// authToken
	VerificationMissing  saerrors.ErrorCode = "verification-missing"
	VerificationConsumed saerrors.ErrorCode = "verification-consumed"
	VerificationExpired  saerrors.ErrorCode = "verification-expired"
	VerificationInvalid  saerrors.ErrorCode = "verification-invalid"
)
