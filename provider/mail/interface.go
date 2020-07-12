package mail

type Client interface {
	SendMail(to []string, from, subject, message string) error
}