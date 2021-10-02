package utils

import (
	"core/errs"
	"bytes"
	"core/logger"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
)

type EmailTemplate struct {
	Message string
	Link    string
	Content string
	Button  string
}

func SendEmail(code string, email, name string) *errs.AppError {
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config : ", err)
	}

	from := config.EmailFrom
	password := config.EmailPassword
	host := "smtp.gmail.com"
	port := config.EmailPort
	toEmail := email
	to := []string{toEmail}
	address := host + ":" + port
	auth := smtp.PlainAuth("", from, password, host)
	subject := ""

	var P EmailTemplate
	var t *template.Template
	t, err = t.ParseFiles("utils/templates/temp.html")
	if err != nil {
		return errs.NewUnexpectedError("error while sending email" + err.Error())
	}
	subject = "Subject:Email verification to signup\n"

	mime := "MIME-version: 1.0;\nContent-Type:text/html; charset=\"UTF-8\";\n"

	if name == "userreg" {
		P.Link = fmt.Sprintf("http://23.99.224.250/v1/auth/verifyemail/%s/%s", email, code)
		P.Message = "Verify Your E-mail Address"
		P.Content = "You're almost ready to get started. Please click on the button below to verify your email address  "
		P.Button = "VERIFY YOUR EMAIL"
	} else if name == "reactivated" {
		subject = "Subject:Email verification to change email\n"
		P.Link = fmt.Sprintf("http://23.99.224.250/v1/auth/verifychangeemail/%s/%s", email, code)
		P.Message = "Verify Your E-mail Address To Change To New Email"
		P.Content = "You are almost ready to get started. Please click on the button below to verify your new email address "
		P.Button = "VERIFY YOUR EMAIL"
	} else {
		subject = "Subject:Password reset link\n"
		P.Link = fmt.Sprintf("http://23.99.224.250/v1/auth/resetpassword/%s/%s", email, code)
		P.Message = "Reset Your Password"
		P.Content = "You're almost ready to get you a new password. Please click on the button below to reset "
		P.Button = "Reset Password"
	}

	buff := new(bytes.Buffer)
	t.Execute(buff, P)
	msg := []byte(subject + mime + buff.String())
	err0 := smtp.SendMail(address, auth, from, to, msg)
	if err0 != nil {
		logger.Error(err0.Error())
		return errs.NewUnexpectedError("error while sending email")
	}
	return nil
}
