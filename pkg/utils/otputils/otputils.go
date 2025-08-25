package otputils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

func GenerateSecureOTP() string {
	max := big.NewInt(1000000)
	n, _ := rand.Int(rand.Reader, max)
	return fmt.Sprintf("%06d", n.Int64())
}

func HashOTP(otp string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(otp))
	return hex.EncodeToString(h.Sum(nil))
}
