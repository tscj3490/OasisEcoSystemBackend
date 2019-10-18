package mail

import (
	"com/merkinsio/oasis-api/config"

	"github.com/sirupsen/logrus"
	gomail "gopkg.in/gomail.v2"
)

/*Request Holds the email send data*/
type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

/*CreateRequest Creates a new Request*/
func CreateRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

/*SendEmail Creates the email dialer and send the email*/
func (r *Request) SendEmail() (bool, error) {
	// smtp dialer
	host := config.Config.GetString("mail.host")
	port := config.Config.GetInt("mail.port")
	username := config.Config.GetString("mail.user")
	password := config.Config.GetString("mail.password")
	d := gomail.NewPlainDialer(host, port, username, password)

	// smtp sender
	s, err := d.Dial()
	if err != nil {
		return false, err
	}

	var from string

	if r.from != "" {
		from = r.from
	} else {
		from = config.Config.GetString("mail.sender")
	}

	// the new message
	m := gomail.NewMessage()

	for _, recipient := range r.to {
		m.SetHeader("From", from)
		m.SetHeader("Subject", r.subject)
		m.SetBody("text/html", r.body)
		m.SetAddressHeader("To", recipient, "")

		if err := gomail.Send(s, m); err != nil {
			logrus.Printf("Could not send email to %q: %v", r.to, err)
		}

		m.Reset()
	}

	return true, nil
}
