package mailer

type Mailer interface {
	Send(to string, subject string, body []byte) error
}
