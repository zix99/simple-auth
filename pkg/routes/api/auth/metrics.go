package auth

import (
	"simple-auth/pkg/instrumentation"
	"simple-auth/pkg/saerrors"
)

var authCounter instrumentation.Counter = instrumentation.NewCounter("sa_auth", "Authentication counter", "type", "success", "errCode")

func incAuthCounterSuccess(name string) {
	authCounter.Inc(name, true, "nil")
}

func incAuthCounterError(name string, err error) {
	authCounter.Inc(name, false, saerrors.UnwrapCode(err))
}
