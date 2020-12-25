package instrumentation

type Counter interface {
	Inc(values ...interface{})
}
