package email

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"
)

func (e *Repository) Send(ctx context.Context, to, templateID string, args map[string]string) error {
	addr := fmt.Sprintf("%s:%d", e.SMTPHost, e.SMTPPort)
	auth := smtp.PlainAuth("", e.Username, e.Password, e.SMTPHost)

	content, err := templateFS.ReadFile("templates/" + templateID + ".tmpl")
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templateID, err.Error())
	}

	tmpl, err := template.New(templateID).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templateID, err.Error())
	}

	var msg strings.Builder
	if err := tmpl.Execute(&msg, args); err != nil {
		return fmt.Errorf("template render error: %w", err)
	}

	return smtp.SendMail(addr, auth, e.From, []string{to}, []byte(msg.String()))
}
