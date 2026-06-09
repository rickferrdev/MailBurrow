package mailer

import (
	"strconv"

	"github.com/rickferrdev/mail-burrow/internal/app/ports"
	"github.com/rickferrdev/mail-burrow/internal/config/env"
	"go.uber.org/fx"
	pkgmail "gopkg.in/mail.v2"
)

var Provide = fx.Provide(New)

type Mailer struct {
	dialer *pkgmail.Dialer
}

func New(env *env.Environment) (ports.Mailer, error) {
	port, err := strconv.Atoi(env.MailerPort)
	if err != nil {
		return nil, err
	}

	dialer := pkgmail.NewDialer(
		env.MailerHost,
		port,
		env.MailerUsername,
		env.MailerPassword,
	)

	mailer := Mailer{
		dialer: dialer,
	}

	return &mailer, nil
}

func (mailer *Mailer) Send(to, from, subject, body string) error {
	message := pkgmail.NewMessage()
	message.SetHeader("To", to)
	message.SetHeader("From", from)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)

	if err := mailer.dialer.DialAndSend(message); err != nil {
		return err
	}

	return nil
}
