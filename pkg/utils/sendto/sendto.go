package sendto

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"

	"github.com/ducklawrence05/go-test-backend-api/config"
)

type EmailAddress struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type Mail struct {
	From    EmailAddress
	To      []string
	Subject string
	Body    string
}

func BuildMessage(mail Mail) string {
	msg := "MIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n"
	msg += fmt.Sprintf("From: %s <%s>\r\n", mail.From.Name, mail.From.Address)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ", "))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)
	return msg
}

func SendTextEmailOtp(smtpCfg *config.SMTP, to []string, otp string) error {
	contentEmail := Mail{
		From:    EmailAddress{Address: smtpCfg.Username, Name: "Duck Test"},
		To:      to,
		Subject: "OTP Verification",
		Body:    fmt.Sprintf("Your OTP is %s. Please enter it to verify your account.", otp),
	}

	messageMail := BuildMessage(contentEmail)

	// send smtp
	auth := smtp.PlainAuth("", smtpCfg.Username, smtpCfg.AppPassword, smtpCfg.Host)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", smtpCfg.Host, smtpCfg.Port),
		auth,
		smtpCfg.Username,
		to,
		[]byte(messageMail),
	)

	if err != nil {
		return err
	}

	return nil
}

func SendTemplateEmailOtp(
	smtpCfg *config.SMTP, to []string,
	nameTemplate string, dataTemplate map[string]any,
) error {
	htmlBody, err := getMailTemplate(nameTemplate, dataTemplate)
	if err != nil {
		return err
	}
	return send(smtpCfg, to, htmlBody)
}

func getMailTemplate(nameTemplate string, dataTemplate map[string]any) (string, error) {
	htmlTemplate := new(bytes.Buffer)
	t := template.Must(template.New(nameTemplate).ParseFiles("templates/email/" + nameTemplate))
	err := t.Execute(htmlTemplate, dataTemplate)
	if err != nil {
		return "", err
	}
	return htmlTemplate.String(), nil
}

func send(smtpCfg *config.SMTP, to []string, htmlTemplate string) error {
	contentEmail := Mail{
		From:    EmailAddress{Address: smtpCfg.Username, Name: "Duck Test"},
		To:      to,
		Subject: "OTP Verification",
		Body:    htmlTemplate,
	}

	messageMail := BuildMessage(contentEmail)

	// send smtp
	auth := smtp.PlainAuth("", smtpCfg.Username, smtpCfg.AppPassword, smtpCfg.Host)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", smtpCfg.Host, smtpCfg.Port),
		auth,
		smtpCfg.Username,
		to,
		[]byte(messageMail),
	)

	if err != nil {
		return err
	}

	return nil
}
