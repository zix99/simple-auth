package engine

import "fmt"

/* stdout engine
Writes emails out to terminal
*/

type stdoutEngine struct{}

func NewStdoutEngine() EmailEngine {
	return &stdoutEngine{}
}

func (s *stdoutEngine) Send(to, from string, data []byte) error {
	fmt.Printf("Email:\nTo: %s\nFrom: %s\n\n%s\n", to, from, string(data))
	return nil
}
