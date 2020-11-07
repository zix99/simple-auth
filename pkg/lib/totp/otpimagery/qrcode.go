package otpimagery

import (
	"simple-auth/pkg/lib/totp"

	qrcode "github.com/skip2/go-qrcode"
)

// GenerateQRCode generates a raw PNG qr code from the totp URI
func GenerateQRCode(otp *totp.Totp, size int) ([]byte, error) {
	return qrcode.Encode(otp.String(), qrcode.Medium, size)
}
