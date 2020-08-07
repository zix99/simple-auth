package otpimagery

import (
	"simple-auth/pkg/lib/totp"

	qrcode "github.com/skip2/go-qrcode"
)

func GenerateQRCode(otp *totp.Totp, size int) ([]byte, error) {
	return qrcode.Encode(otp.String(), qrcode.Medium, size)
}
