package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"

	"email-service/internal/config"
)

type EmailService struct {
	config *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		config: cfg,
	}
}

func (s *EmailService) TestConnection() error {
	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)
	
	conn, err := smtp.Dial(fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort))
	if err != nil {
		return fmt.Errorf("SMTP dial failed: %v", err)
	}
	defer conn.Close()

	if err = conn.StartTLS(&tls.Config{
		ServerName: s.config.SMTPHost,
	}); err != nil {
		return fmt.Errorf("StartTLS failed: %v", err)
	}

	if err = conn.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %v", err)
	}

	return nil
}

func (s *EmailService) SendAuthCode(email, code, authType string) error {
	subject := "Your Authentication Code"

	var body string
	switch authType {
	case "login":
		subject = "Login Verification Code"
		body = fmt.Sprintf(`
			<html>
			<body>
				<h2>Login Verification</h2>
				<p>Your login verification code is: <strong>%s</strong></p>
			</body>
			</html>
		`, code)
	case "registration":
		subject = "Registration Verification Code"
		body = fmt.Sprintf(`
			<html>
			<body>
				<h2>Welcome!</h2>
				<p>Your registration verification code is: <strong>%s</strong></p>
			</body>
			</html>
		`, code)
	default:
		subject = "Your Authentication Code"
		body = fmt.Sprintf(`
			<html>
			<body>
				<h2>Authentication Code</h2>
				<p>Your authentication code is: <strong>%s</strong></p>
			</body>
			</html>
		`, code)
	}

	return s.sendEmail(email, subject, body)
}

func (s *EmailService) sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)

	from := s.config.FromEmail
	if from == "" {
		from = s.config.SMTPUsername
	}

	message := fmt.Sprintf("From: %s\r\n", from) +
		fmt.Sprintf("To: %s\r\n", to) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" +
		body

	conn, err := smtp.Dial(fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort))
	if err != nil {
		return fmt.Errorf("SMTP dial failed: %v", err)
	}
	defer conn.Close()

	if err = conn.StartTLS(&tls.Config{
		ServerName: s.config.SMTPHost,
	}); err != nil {
		return fmt.Errorf("StartTLS failed: %v", err)
	}

	if err = conn.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %v", err)
	}

	if err = conn.Mail(from); err != nil {
		return fmt.Errorf("SMTP mail failed: %v", err)
	}

	if err = conn.Rcpt(to); err != nil {
		return fmt.Errorf("SMTP rcpt failed: %v", err)
	}

	w, err := conn.Data()
	if err != nil {
		return fmt.Errorf("SMTP data failed: %v", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("SMTP write failed: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("SMTP close failed: %v", err)
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}