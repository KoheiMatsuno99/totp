package main

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"totp/ctxtime"

	"github.com/skip2/go-qrcode"
)

type TOTP struct {
	Secret    string
	Digits    int
	Period    int64
	Algorithm string
}

func NewTOTP(secret string) *TOTP {
	return &TOTP{
		Secret:    secret,
		Digits:    6,
		Period:    30,
		Algorithm: "SHA1",
	}
}

func GenerateSecret() (string, error) {
	key := make([]byte, 20)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(key), nil
}

func (t *TOTP) GenerateCode(ctx context.Context) (string, error) {
	ts := ctxtime.Now(ctx).Unix()

	counter := ts / t.Period

	secretBytes, err := base32.StdEncoding.DecodeString(strings.ToUpper(t.Secret))
	if err != nil {
		return "", err
	}

	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(counter))

	mac := hmac.New(sha1.New, secretBytes)
	mac.Write(buf)
	hash := mac.Sum(nil)

	offset := hash[len(hash)-1] & 0x0F
	code := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7FFFFFFF

	otp := fmt.Sprintf("%0*d", t.Digits, code%uint32(math.Pow10(t.Digits)))
	return otp, nil
}

func (t *TOTP) Verify(ctx context.Context, code string) bool {
	baseTs := ctxtime.Now(ctx).Unix()

	for i := -1; i <= 1; i++ {
		ts := baseTs + int64(i)*t.Period
		counter := ts / t.Period

		secretBytes, err := base32.StdEncoding.DecodeString(strings.ToUpper(t.Secret))
		if err != nil {
			return false
		}

		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(counter))

		mac := hmac.New(sha1.New, secretBytes)
		mac.Write(buf)
		hash := mac.Sum(nil)

		offset := hash[len(hash)-1] & 0x0F
		totpCode := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7FFFFFFF

		expectedCode := fmt.Sprintf("%0*d", t.Digits, totpCode%uint32(math.Pow10(t.Digits)))
		if code == expectedCode {
			return true
		}
	}
	return false
}

func (t *TOTP) getQRCodeURL(issuer, account string) string {
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s&digits=%d&period=%d&algorithm=%s",
		issuer, account, t.Secret, issuer, t.Digits, t.Period, t.Algorithm)
}

func (t *TOTP) GenerateQRCode(issuer, account, filename string) error {
	url := t.getQRCodeURL(issuer, account)
	return qrcode.WriteFile(url, qrcode.Medium, 256, filename)
}
