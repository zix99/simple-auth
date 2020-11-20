package engine

type EmailEngine interface {
	Send(to, from string, data []byte) error
}
