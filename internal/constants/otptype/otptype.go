package otptype

type OTPType string

const (
	OTPRegister       OTPType = "register"
	OTPForgotPassword OTPType = "forgot_password"
	OTPChangeEmail    OTPType = "change_email"
	OTPRestoreAccount OTPType = "restore_account"
)
