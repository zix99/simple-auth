package config

var (
	baseTrue  = true
	baseFalse = false
	TruePtr   = &baseTrue
	FalsePtr  = &baseFalse
)

func IntPtr(v int) *int {
	return &v
}

func BoolPtr(v bool) *bool {
	return &v
}

func StrPtr(s string) *string {
	return &s
}

func CoalesceInt(a *int, b *int) *int {
	if a != nil {
		return a
	}
	if b != nil {
		return b
	}
	panic("coalesce")
}

func CoalesceBool(a *bool, b *bool) *bool {
	if a != nil {
		return a
	}
	if b != nil {
		return b
	}
	panic("coalesce")
}

func CoalesceString(a *string, b *string) *string {
	if a != nil {
		return a
	}
	if b != nil {
		return b
	}
	panic("coalesce")
}
