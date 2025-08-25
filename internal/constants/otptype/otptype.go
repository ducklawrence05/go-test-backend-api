package otptype

type OTPType string

const (
	OTPEmailVerify    OTPType = "email_verify"
	OTPForgotPassword OTPType = "forgot_password"
	OTPChangeEmail    OTPType = "change_email"
)
