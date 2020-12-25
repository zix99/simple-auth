// +build !prometheus

package instrumentation

type mockCounter struct{}

func (s *mockCounter) Inc(values ...interface{}) {}

func NewCounter(name, help string, labels ...string) Counter {
	return &mockCounter{}
}
