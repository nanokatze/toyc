package compiler

// TODO: make these nicer
// TODO: use some prefix other than "Build"?

func BuildMakeTuple(b *Builder, elems ...*Class) *Class {
	typ := make([]Type, len(elems))
	for i, v := range elems {
		typ[i] = v.Type()
	}
	return b.Value2(OpMakeTuple, MakeTupleType(typ...), nil, elems...)
}

func BuildTupleExtract(b *Builder, tup *Class, idx int) *Class {
	return b.Value2(OpTupleExtract, tup.Type().(*TupleType).Elems()[idx], idx, tup)
}

func BuildConst(b *Builder, typ Type, imm int64) *Class {
	return b.Value2(OpConst, typ, imm)
}
