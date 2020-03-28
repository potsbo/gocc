package types

type Kind int

const (
	_ Kind = iota
	Int
	Pointer
)

type Type interface {
	Kind() Kind
	PointingTo() Type
}

func (k Kind) Size() int {
	switch k {
	case Int:
		return 4 // TODO
	case Pointer:
		return 4
	}
	return 0
}

type typeImpl struct {
	kind       Kind
	pointingTo Type
}

func (t typeImpl) Kind() Kind {
	return t.kind
}
func (t typeImpl) PointingTo() Type {
	return t.pointingTo
}

func (k Kind) Identifier() string {
	switch k {
	case Int:
		return "int"
	}

	panic("Unreachable code")
}

func (k Kind) Type() Type {
	return &typeImpl{kind: k}
}

func PointingTo(t Type) Type {
	return &typeImpl{kind: Pointer, pointingTo: t}
}

func All() []Kind {
	return []Kind{
		Int,
	}
}

func NewInt() Type {
	return typeImpl{kind: Int}
}
