package otptype

type OTPType string

const (
	Register       OTPType = "register"
	ForgotPassword OTPType = "forgot_password"
	ChangeEmail    OTPType = "change_email"
	RestoreAccount OTPType = "restore_account"
)
