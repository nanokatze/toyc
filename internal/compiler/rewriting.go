package compiler

import (
	"fmt"
)

// TODO: rename?
type RewriteContext struct {
	b     *Builder
	seen  map[*Value]bool
	stack []*Value
}

func (r *RewriteContext) B() *Builder { return r.b }

// TODO: methods for killing the matched value? I guess actually we could make
// seen a map[*Value]bool and set it to false for killed values.

// TODO: rename, e.g. Value
func (r *RewriteContext) Add2(op Op, typ Type, imm any, args ...*Class) {
	r.add(r.b.Sea.value(op, typ, imm, args...))
}

// TODO: rename?
// TODO: rewrite this to be "safe" (always go through Add2)? Or document the
// requirements that Class must be from the same Sea.
func (r *RewriteContext) Class(c *Class) {
	for _, v := range c.values {
		r.add(v)
	}
}

func (r *RewriteContext) add(v *Value) {
	if _, ok := r.seen[v]; ok {
		return // we saw this value already
	}
	r.seen[v] = true
	r.stack = append(r.stack, v)
}

// TODO: rename?
func (rc *RewriteContext) applyRules() {
	for len(rc.stack) > 0 {
		v := rc.stack[len(rc.stack)-1]
		rc.stack = rc.stack[:len(rc.stack)-1]

		for _, rule := range rc.b.Rules {
			if rule.Pattern.Match(v) {
				rule.Rewrite(rc, v)
			}
		}
	}
}

type RewriteRule struct {
	Name    string
	Pattern *Pattern
	Rewrite func(*RewriteContext, *Value) // TODO: I guess we could also eliminate *Value and pass that through *RewriteContext
}

// This handles only 2-ary ops, TODO: teach it to handle n-ary by only swapping
// the first two args?
// TODO: rename to something like CommutativityRule or whatever.
func Commutativity(op Op) RewriteRule {
	return RewriteRule{
		Name: fmt.Sprintf("%v commutativity", op),
		Pattern: &Pattern{
			Op:   op,
			Args: []*Pattern{{}, {}},
		},
		Rewrite: func(rc *RewriteContext, v *Value) {
			rc.Add2(v.Op(), v.Type(), v.Imm(), v.Arg(1), v.Arg(0))
		},
	}
}

/*
func Associativity(op Op) RewriteRule {
	return RewriteRule{
		Name: fmt.Sprintf("%v associativity", op),
		Pattern: &Pattern{
			Op: op,
			Args: []*Pattern{
				{},
				{},
			},
		},
		Replace: func(b *Builder, rc *RewriteContext, v *Value) {
		},
	}
}
*/

type Builder struct {
	Sea   *Sea
	Rules []RewriteRule
}

// TODO: rename to Build? just Value?
func (b *Builder) Value2(op Op, typ Type, imm any, args ...*Class) *Class {
	// TODO: reuse this with a sync.Pool
	rc := &RewriteContext{
		b:    b,
		seen: make(map[*Value]bool),
	}

	rc.Add2(op, typ, imm, args...)

	rc.applyRules()

	return rc.b.Sea.class(typ, rc.seen)
}

func Rewrite(b *Builder, c *Class) *Class {
	// TODO: see if we can rewrite this non-recursively

	visited := make(map[*Class]*Class)

	var f func(c *Class) *Class
	f = func(c *Class) *Class {
		if x, ok := visited[c]; ok {
			return x
		}

		rc := &RewriteContext{
			b:    b,
			seen: make(map[*Value]bool),
		}
		for v := range c.Values() {
			// TODO: factor out remapping the args. This also would come in
			// useful in Sea.value.
			args := make([]*Class, len(v.Args()))
			for i := range args {
				args[i] = f(v.Arg(i))
			}
			// Go through value creation path, as we don't know whether it's
			// from the same sea or not.
			rc.Add2(v.Op(), v.Type(), v.Imm(), args...)
		}
		rc.applyRules()

		// TODO: factor this out into a method on RewriteContext, possibly fold
		// into applyRules
		x := rc.b.Sea.class(c.Type(), rc.seen)

		visited[c] = x
		return x
	}
	return f(c)
}
