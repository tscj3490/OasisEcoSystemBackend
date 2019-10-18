package mail

/*Main Main email template data*/
type Main struct {
	Content string
}

/*SendRegistrationPassword Send registration password email template data*/
type SendRegistrationPassword struct {
	Name     string
	Password string
}

type SendRecoveryPassword struct {
	Token string
}
