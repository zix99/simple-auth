package appcontext

type Context interface {
	Get(key string) interface{}
}

type RWContext interface {
	Context
	Set(key string, i interface{})
}

type ProviderFunc func(ctx Context) (key string, i interface{})

func With(key string, i interface{}) ProviderFunc {
	return func(ctx Context) (key string, i interface{}) {
		return key, i
	}
}

type Container struct {
	items map[string]interface{}
}

func NewContainer() *Container {
	return &Container{
		items: make(map[string]interface{}),
	}
}

func (s *Container) Get(key string) interface{} {
	return s.items[key]
}

func (s *Container) Set(key string, i interface{}) {
	s.items[key] = i
}

func (s *Container) Use(f ProviderFunc) {
	if key, val := f(s); val != nil {
		s.items[key] = val
	}
}
