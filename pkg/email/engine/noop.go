package engine

import "github.com/sirupsen/logrus"

/* noop engine
Does nothing when receiving emails
*/

type noopEngine struct {
	returns error
}

func NewNoopEngine(returns error) EmailEngine {
	return &noopEngine{
		returns,
	}
}

func (s *noopEngine) Send(to, from string, data []byte) error {
	logrus.Warnf("Email to %s is going into the void!", to)
	return s.returns
}
