package core

import (
	"fmt"

	"toyc/internal/compiler"
)

// TODO: introduce an op to convert between FloatType and Int
type FloatType struct {
	E int
	M int
}

func (t FloatType) String() string { return fmt.Sprintf("Float[%d,%d]", t.E, t.M) }

var FloatE8M23 = FloatType{8, 23}

var (
	OpFAdd = compiler.DefOp("FAdd", nil)
	OpFSub = compiler.DefOp("FSub", nil)
	OpFMul = compiler.DefOp("FMul", nil)
	OpFDiv = compiler.DefOp("FDiv", nil)

	OpFMin = compiler.DefOp("FMin", nil)
	OpFMax = compiler.DefOp("FMax", nil)

	OpFFloor = compiler.DefOp("FFloor", nil)
	OpFCeil  = compiler.DefOp("FCeil", nil)

	OpFEqual       = compiler.DefOp("FEqual", nil)
	OpFLess        = compiler.DefOp("FLess", nil)
	OpFLessOrEqual = compiler.DefOp("FLessOrEqual", nil)
)

func init() {
	Rules = append(Rules,
		compiler.Commutativity(OpFAdd),
		compiler.Commutativity(OpFMul),

		compiler.Commutativity(OpFMin),
		compiler.Commutativity(OpFMax),

		compiler.Commutativity(OpFEqual),
	)
}
