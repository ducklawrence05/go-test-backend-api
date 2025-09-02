package otputils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateSecureOTP() string {
	max := big.NewInt(1000000)
	n, _ := rand.Int(rand.Reader, max)
	return fmt.Sprintf("%06d", n.Int64())
}
