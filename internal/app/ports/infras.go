package ports

type Mailer interface {
	Send(to, from, subject, body string) error
}
