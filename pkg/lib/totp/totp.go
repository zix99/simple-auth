package totp

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type Totp struct {
	secret  []byte
	Subject string
	Issuer  string
}

func CreateSecret(keylen int) ([]byte, error) {
	key := make([]byte, keylen)
	if n, err := rand.Read(key); err != nil || n != keylen {
		return nil, err
	}
	return key, nil
}

func EncodeSecretb32(secret []byte) string {
	return strings.ToUpper(base32.StdEncoding.EncodeToString(secret))
}

func NewTOTP(keylen int, issuer, subject string) (*Totp, error) {
	key, err := CreateSecret(keylen)
	if err != nil {
		return nil, err
	}

	return &Totp{
		secret:  key,
		Subject: subject,
		Issuer:  issuer,
	}, nil
}

func FromSecret(b32 string, issuer, subject string) (*Totp, error) {
	key, err := base32.StdEncoding.DecodeString(b32)
	if err != nil {
		return nil, err
	}
	return &Totp{
		secret:  key,
		Subject: subject,
		Issuer:  issuer,
	}, nil
}

func FromURL(url *url.URL) (*Totp, error) {
	query := url.Query()

	if query.Get("digits") != "6" {
		return nil, errors.New("invalid digits")
	}
	if query.Get("period") != "30" {
		return nil, errors.New("invalid period")
	}
	if query.Get("algorithm") != "SHA1" {
		return nil, errors.New("unsupported algorithm")
	}

	key, err := base32.StdEncoding.DecodeString(query.Get("secret"))
	if err != nil {
		return nil, err
	}

	parts := strings.Split(url.Path[1:], ":")
	if len(parts) != 2 {
		return nil, errors.New("invalid subject")
	}

	return &Totp{
		secret:  key,
		Issuer:  parts[0],
		Subject: parts[1],
	}, nil
}

func ParseTOTP(uri string) (*Totp, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	return FromURL(parsed)
}

func (s *Totp) URL() *url.URL {
	return &url.URL{
		Scheme: "otpauth",
		Host:   "totp",
		Path:   s.Issuer + ":" + s.Subject,
		RawQuery: url.Values{
			"secret":    {s.Secret()},
			"issuer":    {s.Issuer},
			"algorithm": {"SHA1"},
			"digits":    {"6"},
			"period":    {"30"},
		}.Encode(),
	}
}

func (s *Totp) String() string {
	return s.URL().String()
}

func (s *Totp) Secret() string {
	return EncodeSecretb32(s.secret)
}

func (s *Totp) GetHOTP(interval int64) string {
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, uint64(interval))

	hash := hmac.New(sha1.New, s.secret)
	hash.Write(bs)
	h := hash.Sum(nil)

	// read number
	o := (h[19] & 0xf)
	header := binary.BigEndian.Uint32(h[o : o+4])
	h12 := (int(header) & 0x7fffffff) % 1000000

	return fmt.Sprintf("%06d", h12)
}

func (s *Totp) GetTOTP() string {
	interval := time.Now().Unix() / 30
	return s.GetHOTP(interval)
}

// Validate the code, allowing a certain amount of time-drift
func (s *Totp) Validate(code string, drift int) bool {
	interval := time.Now().Unix() / 30
	if code == s.GetHOTP(interval) {
		return true
	}

	for i := 1; i < drift; i++ {
		if code == s.GetHOTP(interval+int64(i)) {
			return true
		}
		if code == s.GetHOTP(interval-int64(i)) {
			return true
		}
	}
	return false
}
