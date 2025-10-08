package compiler

// TODO: make the fields private and introduce accessors
type Value struct {
	Op   *Op
	Type Type
	Args []*EqValues
	Imm  any
}

// TODO: make this a struct so that the user can't fuck around with the
// internals, and also eq values should have a type attached to them.
// TODO: we might want to nih an orderedset so that we can iterate in the
// insertion order and have reproducible prints and everything.
type EqValues map[*Value]struct{}

func (set *EqValues) Type() Type {
	for v := range *set {
		return v.Type
	}
	panic("unreachable")
}

func (set *EqValues) Add(v *Value) {
	(*set)[v] = struct{}{}
}

// TODO: make this a standalone function?
func (a *EqValues) Merge(b *EqValues) {
	for v := range *b {
		a.Add(v)
	}
	*b = *a
}
