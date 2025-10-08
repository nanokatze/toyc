package core

import (
	"fmt"

	"toyc/internal/compiler"
)

// TODO: bring ArrayType into a usable state

type ArrayType struct {
	len  int64
	elem compiler.Type
}

// TODO: hash cons this stuff later
func MakeArrayType(len int64, elem compiler.Type) ArrayType {
	return ArrayType{len, elem}
}

func (typ ArrayType) Len() int64 { return typ.len }

func (typ ArrayType) Elem() compiler.Type { return typ.elem }

func (typ ArrayType) String() string { return fmt.Sprintf("[%d]%v", typ.len, typ.elem) }

var OpMakeArray = compiler.DefOp("MakeArray", validateMakeArray)

// TODO: make this take elem type explicitly?
func MakeArray(b *compiler.Builder, et compiler.Type, args ...*compiler.Class) *compiler.Class {
	t := MakeArrayType(int64(len(args)), et)
	return b.Value2(OpMakeArray, t, nil, args...)
}

var OpArrayExtract = compiler.DefOp("ArrayExtract", validateArrayExtract)

func ArrayExtract(b *compiler.Builder, arr *compiler.Class, idx int64) *compiler.Class {
	return b.Value2(OpArrayExtract, arr.Type().(ArrayType).Elem(), idx, arr)
}

func init() {
	Rules = append(Rules,
		compiler.RewriteRule{
			Name: "Forward ArrayExtract to the element's definition",
			Pattern: &compiler.Pattern{
				Op:   OpArrayExtract,
				Args: []*compiler.Pattern{{Op: OpMakeArray, ArgsDDD: true}},
			},
			Rewrite: func(rc *compiler.RewriteContext, v *compiler.Value) {
				idx := v.Imm().(int64)
				for arr := range v.Arg(0).Values() {
					if arr.Op() == OpMakeArray {
						rc.Class(arr.Arg(int(idx)))
					}
				}
			},
		},
		compiler.RewriteRule{
			Name: "Split CondSelect of arrays into per element CondSelect", // TODO: better name?
			Pattern: &compiler.Pattern{
				Op:   OpCondSelect,
				Args: []*compiler.Pattern{{}, {}, {}},
			},
			Rewrite: func(rc *compiler.RewriteContext, v *compiler.Value) {
				arr, ok := v.Type().(ArrayType)
				if !ok {
					return
				}

				x := v.Arg(0)
				y := v.Arg(1)
				cond := v.Arg(2)

				elems := make([]*compiler.Class, arr.Len())
				for i := range arr.Len() {
					x_i := ArrayExtract(rc.B(), x, i)
					y_i := ArrayExtract(rc.B(), y, i)
					elems[i] = CondSelect(rc.B(), x_i, y_i, cond)
				}
				rc.Add2(OpMakeArray, arr, nil, elems...)
			},
		})
}

func validateMakeArray(typ compiler.Type, imm any, args ...*compiler.Class) {
	if imm != nil {
		panic("imm must be nil")
	}

	arr := typ.(ArrayType)
	if len(args) != int(arr.Len()) {
		panic("mismatched len")
	}
	elem := arr.Elem()
	for _, a := range args {
		if a.Type() != elem {
			panic("mismatched types")
		}
	}
}

func validateArrayExtract(typ compiler.Type, imm any, args ...*compiler.Class) {
	if len(args) != 1 {
		panic("wrong args")
	}

	idx := imm.(int64)

	arr := args[0]

	arrType := arr.Type().(ArrayType)

	if arrType.Elem() != typ {
		panic("type mismatch")
	}
	if !(0 <= idx && idx < arrType.len) {
		panic("index out of bounds")
	}
}
