package compiler

import (
	"slices"
	"strconv"
)

type Type any

type TupleType struct{ elems []Type }

// TODO: hash cons tuple types
func MakeTupleType(elems []Type) *TupleType {
	return &TupleType{elems: slices.Clone(elems)}
}

func (t *TupleType) String() string {
	// TODO: iterate over elems
	return "()"
}

func (t *TupleType) Elems() []Type {
	return t.elems
}

type BitsType struct{ n int64 }

// TODO: hash cons bits types
func MakeBitsType(n int64) *BitsType {
	return &BitsType{n: n}
}

func (t *BitsType) String() string { return "Bits(" + strconv.FormatInt(t.n, 10) + ")" }

var (
	TypeBits1  = MakeBitsType(1)
	TypeBits8  = MakeBitsType(8)
	TypeBits16 = MakeBitsType(16)
	TypeBits32 = MakeBitsType(32)
)
