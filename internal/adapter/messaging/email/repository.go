package email

import (
	"crypto/tls"
	"embed"
	"fmt"
	"labra/internal/entity"
	"net/smtp"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

type Repository struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	From     string
}

func NewRepository(smtpHost string, smtpPort int, username, password, from string) *Repository {
	return &Repository{
		SMTPHost: smtpHost,
		SMTPPort: smtpPort,
		Username: username,
		Password: password,
		From:     from,
	}
}

func (e *Repository) Ping() error {
	addr := fmt.Sprintf("%s:%d", e.SMTPHost, e.SMTPPort)
	c, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer c.Close()

	// Upgrade to TLS if supported
	if ok, _ := c.Extension("STARTTLS"); ok {
		tlsconfig := &tls.Config{
			ServerName: e.SMTPHost,
		}
		if err = c.StartTLS(tlsconfig); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	}

	auth := smtp.PlainAuth("", e.Username, e.Password, e.SMTPHost)
	if err := c.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	return nil
}

func (e *Repository) GetType() entity.ContactType {
	return entity.ContactTypeEmail
}
