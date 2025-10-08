package compiler

var (
	OpMakeTuple    = DefOp("MakeTuple", nil)
	OpTupleExtract = DefOp("TupleExtract", nil)
)

var OpConst = DefOp("Const", validateConst)

func validateConst(typ Type, imm any, args ...*Class) {
	_ = imm.(int64)
}

var OpEqual = DefOp("Equal", nil)

var OpCondSelect = DefOp("CondSelect", nil)

// OpIAdd = DefOp("IAdd")
// OpISub = DefOp("ISub")
// OpIMul = DefOp("IMul")
