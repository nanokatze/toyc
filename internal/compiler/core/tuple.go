package core

/*
type TupleType struct {
	// TODO: make it private and provide a Elems() method instead?
	elems []Type
}

// TODO: hash cons tuple types
// TODO: make this accept an iterator or idk because certain users have to write
// annoying loop that would be better served by writing basically a map.
func MakeTupleType(elems ...Type) *TupleType {
	return &TupleType{slices.Clone(elems)}
}

func (t *TupleType) String() string {
	elems := make([]string, len(t.elems))
	for i, e := range t.elems {
		elems[i] = e.String()
	}

	return "(" + strings.Join(elems, ", ") + ")"
}

func (t *TupleType) Elems() []Type {
	return t.elems
}

var (
	OpMakeTuple    = DefOp("MakeTuple", nil)
	OpTupleExtract = DefOp("TupleExtract", nil)
)

func MakeTuple(b *Builder, elems ...*Class) *Class {
	typ := make([]Type, len(elems))
	for i, v := range elems {
		typ[i] = v.Type()
	}
	return b.Value2(OpMakeTuple, MakeTupleType(typ...), nil, elems...)
}

func TupleExtract(b *Builder, tup *Class, idx int) *Class {
	return b.Value2(OpTupleExtract, tup.Type().(*TupleType).Elems()[idx], idx, tup)
}
*/
