package compiler

import (
	"slices"
)

// TODO: move CSE into its own file?

// TODO: contemplate removing at least some indirections

// TODO: we need a special hash and equality on this so that we can handle
// commutative ops, hash would combine args[0] and args[1] with a commutative
// operator like xor
type cseKey struct {
	op   *Op
	typ  Type
	args [3]*EqValues // aaaaaaaaaa
	imm  any
}

// TODO: give this object a better name
type CSE struct {
	// TODO: use a custom hashmap so we can use a slice of EqValues
	mv map[cseKey]*Value

	mc map[*Value]*EqValues

	// TODO: we probably also want a *EqValues -> *EqValues map for when we
	// merge equality classes?
}

func NewCSE() *CSE {
	return &CSE{
		mv: make(map[cseKey]*Value),
		mc: make(map[*Value]*EqValues),
	}
}

func (cse *CSE) Value(op *Op, typ Type, args []*EqValues, imm any) *Value {
	var args2 [3]*EqValues // ugh
	if copy(args2[:], args) != len(args) {
		panic("aaaaaaaa")
	}

	k := cseKey{op, typ, args2, imm}

	v, ok := cse.mv[k]
	if !ok {
		v = &Value{op, typ, slices.Clone(args), imm}
		cse.mv[k] = v
	}
	return v
}

/*
func (cse *CSE) Class(v *Value) *EqValues {
	eqc, ok := cse.mc[v]
	if !ok {
		eqc = &EqValues{}
		eqc.Add(v)
		cse.mc[v] = eqc
	}
	return eqc
}
*/

func (cse *CSE) DeclareEquivalent(vs ...*Value) *EqValues {
	var eqc *EqValues
	// TODO: choose eqc in a better way
	for _, v := range vs {
		eqc2, ok := cse.mc[v]
		if ok {
			eqc = eqc2
			break
		}
	}
	if eqc == nil {
		eqc = &EqValues{}
	}

	for _, v := range vs {
		eqc2, ok := cse.mc[v]
		if ok {
			eqc.Merge(eqc2)
		}
		eqc.Add(v)
		cse.mc[v] = eqc
	}

	return eqc
}

type Builder struct {
	CSE          *CSE
	RewriteRules []RewriteRule
}

func (b Builder) Value(op *Op, typ Type, args []*EqValues, imm any) *EqValues {
	v := b.CSE.Value(op, typ, args, imm)

	vs := []*Value{v}    // set of equivalent values; TODO: make this actually a set
	stack := []*Value{v} // stack of values that need to be rewritten
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, rule := range b.RewriteRules {
			if binds, ok := rule.Match(v); ok {
				// TODO: we need to apply rewrites recursively. We should keep a
				// queue of *Values we need to apply rewrites to and when we get a
				// new *Value, we need try to apply rewrite.
				u := rule.Replace(b, binds...)
				if u == v {
					continue
				}

				// ehh
				if slices.Index(vs, u) == -1 {
					vs = append(vs, u)
					stack = append(stack, u)
				}
			}
		}
	}

	return b.CSE.DeclareEquivalent(vs...)
}
