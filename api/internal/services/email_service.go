package services

import (
	"bytes"
	"embed"
	"embox/internal/config"
	"html/template"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	config *config.EmailConfig
}

type MailData struct {
	Language string
	Subject  string
	Body     template.HTML
}

//go:embed email/default.html
var emailTemplateFS embed.FS

func NewEmailService(config *config.EmailConfig) *EmailService {
	return &EmailService{config: config}
}

func (s *EmailService) SendEmail(ctxLang any, to, subject, body string) error {
	htmlBody, err := renderMailTemplate(MailData{
		Language: getLanguage(ctxLang),
		Subject:  subject,
		Body:     template.HTML(body),
	})
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(s.config.Host, s.config.Port, s.config.From, s.config.Password)

	return d.DialAndSend(m)
}

func renderMailTemplate(data MailData) (string, error) {
	tmpl, err := template.ParseFS(emailTemplateFS, "email/default.html")
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func getLanguage(ctxLang any) string {
	lang, ok := ctxLang.(string)
	if !ok {
		return "de"
	}
	if len(lang) > 2 {
		return lang[:2]
	}
	return lang
}
