package engine

type MockEmailEngine struct {
	returns error
	sends   int
	last    string
}

func NewMockEngine(returns error) *MockEmailEngine {
	return &MockEmailEngine{returns, 0, ""}
}

func (s *MockEmailEngine) Send(to, from string, data []byte) error {
	s.last = string(data)
	s.sends++
	return s.returns
}

func (s *MockEmailEngine) SendCount() int {
	return s.sends
}

func (s *MockEmailEngine) LastEmail() string {
	return s.last
}

func (s *MockEmailEngine) Reset() {
	s.sends = 0
	s.last = ""
}
