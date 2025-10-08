package main

import (
	"toyc/internal/compiler"
	"toyc/internal/compiler/core"
)

func main() {
	sea := compiler.NewSea()
	b := compiler.Builder{Sea: sea}

	one := b.Value2(core.OpIConst, core.Int32, int64(1))
	two := b.Value2(core.OpIConst, core.Int32, int64(2))
	or := b.Value2(core.OpOr, core.Int32, nil, one, two)

	compiler.Dump(sea, or, nil)
}
