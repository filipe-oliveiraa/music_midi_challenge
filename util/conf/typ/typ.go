package typ

type Type interface {
	GetValue(v string) (any, error)
}
