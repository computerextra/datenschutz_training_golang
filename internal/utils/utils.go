package utils

import (
	"fmt"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	gomail "gopkg.in/gomail.v2"
)

func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func SendVerificationMail(to, token, uri string) error {
	from, ok := os.LookupEnv("EMAIL_FROM")
	if !ok {
		return fmt.Errorf("env error: EMAIL_FROM missing")
	}
	server, ok := os.LookupEnv("EMAIL_SERVER")
	if !ok {
		return fmt.Errorf("env error: EMAIL_SERVER missing")
	}
	portStr, ok := os.LookupEnv("EMAIL_PORT")
	if !ok {
		return fmt.Errorf("env error: EMAIL_PORT missing")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}
	username, ok := os.LookupEnv("MAIL_USER")
	if !ok {
		return fmt.Errorf("env error: MAIL_USER missing")
	}
	pass, ok := os.LookupEnv("EMAIL_PASSWORD")
	if !ok {
		return fmt.Errorf("env error: EMAIL_PASSWORD missing")
	}

	link := fmt.Sprintf("%s/verify?token=%s", uri, token)

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Bitte bestätige deine E-Mail-Adresse")
	m.SetBody("text/html", "Bitte klicke auf den folenden Link, um deine E-Mail-Adresse zu bestätigen: <a href=\""+link+"\">Bestätigen</a>")

	d := gomail.NewDialer(server, port, username, pass)
	return d.DialAndSend(m)
}
