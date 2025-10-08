package main

import (
	"log"

	"toy/internal/compiler"
)

// TODO: make a dsl for rewrite rules or something
var rules = []compiler.RewriteRule{
	/*
		// TODO: this rule is actually a bit tricky. if we have several
		// OpMakeTuples in the same eq class, that means the argument passes to
		// those OpMakeTuples should also live in their respective eq classes.
		{
			Name: "TupleExtract from TupleMake",
			Pattern: &compiler.Pattern{
				Op: compiler.OpTupleExtract,
				Args: []*compiler.Pattern{
					{
						Op:      compiler.OpTupleMake,
						ArgsDDD: true,
						Bind:    -1,
					},
				},
				Bind: 0,
			},
			Binds: 1,
			Replace: func(b compiler.Builder, binds ...*compiler.Value) *compiler.Value {
				v := binds[0]
				// v.Imm.(int)
			},
		},
	*/
	{
		Name: "constant fold IAdd",
		Pattern: &compiler.Pattern{
			Op: compiler.OpIAdd,
			Args: []*compiler.Pattern{
				{
					Op:   compiler.OpConst,
					Bind: 0,
				},
				{
					Op:   compiler.OpConst,
					Bind: 1,
				},
			},
			Bind: 2,
		},
		Binds: 3,
		Replace: func(b compiler.Builder, binds ...*compiler.Value) *compiler.Value {
			v := binds[2]
			x := binds[0]
			y := binds[1]
			return b.CSE.Value(compiler.OpConst, v.Type, nil, x.Imm.(int64)+y.Imm.(int64))
		},
	},
	{
		Name: "constant fold IMul",
		Pattern: &compiler.Pattern{
			Op: compiler.OpIMul,
			Args: []*compiler.Pattern{
				{
					Op:   compiler.OpConst,
					Bind: 0,
				},
				{
					Op:   compiler.OpConst,
					Bind: 1,
				},
			},
			Bind: 2,
		},
		Binds: 3,
		Replace: func(b compiler.Builder, binds ...*compiler.Value) *compiler.Value {
			v := binds[2]
			x := binds[0]
			y := binds[1]
			return b.CSE.Value(compiler.OpConst, v.Type, nil, x.Imm.(int64)*y.Imm.(int64))
		},
	},
}

func main() {
	b := compiler.Builder{
		CSE:          compiler.NewCSE(),
		RewriteRules: rules,
	}

	one := b.Value(compiler.OpConst, compiler.TypeBits32, nil, int64(1))
	two := b.Value(compiler.OpConst, compiler.TypeBits32, nil, int64(2))
	sum := b.Value(compiler.OpIAdd, compiler.TypeBits32, []*compiler.EqValues{one, two}, nil)
	// sumOfSums := b.Value(compiler.OpIAdd, compiler.TypeBits32, []*compiler.EqValues{sum, sum}, nil)
	// sumTimesTwo := b.Value(compiler.OpIMul, compiler.TypeBits32, []*compiler.EqValues{sum, two}, nil)

	for v := range *sum {
		// log.Println("-----")
		log.Println(v)
		// log.Println("begin args")
		// for _, arg := range v.Args {
		// 	for a := range *arg {
		// 		log.Println(a)
		// 		break
		// 	}
		// }
		// log.Println("end args")
	}
}
