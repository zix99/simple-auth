package engine

type noopEngine struct {
	returns error
}

func NewNoopEngine(returns error) EmailEngine {
	return &noopEngine{
		returns,
	}
}

func (s *noopEngine) Send(to, from string, data []byte) error {
	return s.returns
}
