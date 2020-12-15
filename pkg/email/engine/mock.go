package engine

import (
	"sync/atomic"
)

/* Mock Engine
Used primarily for unit-testing
*/

type MockEmailEngine struct {
	returns error
	sends   int32
	last    string
}

func NewMockEngine(returns error) *MockEmailEngine {
	return &MockEmailEngine{returns, 0, ""}
}

func (s *MockEmailEngine) Send(to, from string, data []byte) error {
	s.last = string(data)
	atomic.AddInt32(&s.sends, 1)
	return s.returns
}

func (s *MockEmailEngine) SendCount() int {
	return int(atomic.LoadInt32(&s.sends))
}

func (s *MockEmailEngine) LastEmail() string {
	return s.last
}

func (s *MockEmailEngine) Reset() {
	s.sends = 0
	s.last = ""
}
