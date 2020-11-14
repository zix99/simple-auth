package ui

import "simple-auth/pkg/saerrors"

const (
	// Common
	errorInvalidAccount saerrors.ErrorCode = "invalid-account"
	errorEmailSend      saerrors.ErrorCode = "email-send"

	// Create account
	errorInvalidRecaptcha saerrors.ErrorCode = "invalid-recaptcha"
)
