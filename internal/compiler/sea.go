package compiler

import (
	"cmp"
	"iter"
	"maps"
	"slices"
	"strconv"
)

type Type interface {
	// TODO: more constraints
	String() string
}

type ValueID int32

func (id ValueID) String() string { return "v" + strconv.FormatInt(int64(id), 10) }

// Value defines how a value is computed.
type Value struct {
	id    ValueID
	op    Op
	typ   Type
	imm   any
	args  []*Class // TODO: small storage for Args?
	class *Class
}

func (v *Value) ID() ValueID { return v.id }

func (v *Value) Op() Op { return v.op }

func (v *Value) Type() Type { return v.typ }

func (v *Value) Imm() any { return v.imm }

func (v *Value) Args() []*Class { return v.args }

func (v *Value) Arg(i int) *Class { return v.args[i] }

// Class that this value is in. May be nil while applying rewrite rules.
func (v *Value) Class() *Class { return v.class }

type ClassID int32

func (id ClassID) String() string { return "c" + strconv.FormatInt(int64(id), 10) }

// Class is a set of equivalent value definitions.
type Class struct {
	id      ClassID
	typ     Type
	classes []*Class // TODO: make small storage for these?
	values  []*Value
	users   map[*Value]struct{} // may be mutated; TODO: fix iteration order for this somehow
	merged  *Class              // may be mutated
}

// TODO: rename
func (c *Class) Newest() *Class {
	// TODO: merged should always point to the newest class, i.e. if c.merged is
	// not nil, then c.merged.merged must always be nil
	for c.merged != nil {
		c = c.merged
	}
	return c
}

func (c *Class) ID() ClassID { return c.id }

func (c *Class) Type() Type { return c.typ }

func (c *Class) Values() iter.Seq[*Value] {
	return func(yield func(*Value) bool) {
		for _, c := range c.classes {
			for v := range c.Values() {
				if !yield(v) {
					return
				}
			}
		}

		for _, v := range c.values {
			if !yield(v) {
				return
			}
		}
	}
}

// TODO: rename to something else to make it clearer what the requirements
func (c *Class) Value() *Value {
	// TODO: assert that there's exactly one value in this class.
	for v := range c.Values() {
		return v
	}
	panic("unreachable")
}

// Users returns an iterator over the values that use this equivalence class.
//
// TODO: do we need a method to test if a particular value uses this class?
func (c *Class) Users() iter.Seq[*Value] {
	return maps.Keys(c.users)
}

// TODO: rename
type seaKey struct {
	op   Op
	typ  Type
	imm  any
	args [10]*Class // aaaaaaaaa
}

func mkseakey(op Op, typ Type, imm any, args ...*Class) seaKey {
	var args_ [10]*Class
	for i, a := range args {
		args_[i] = a
	}

	k := seaKey{op, typ, imm, args_}

	return k
}

type Sea struct {
	m map[seaKey]*Value

	vid int32
	cid int32

	// Debugging

	values  map[ValueID]*Value
	classes map[ClassID]*Class
}

// TODO: accept flags for e.g. disabling debugging for perf
func NewSea() *Sea {
	return &Sea{
		m: make(map[seaKey]*Value),

		// TODO: only create these if we pass a certain flag to NewSea
		// (validation/debug.)
		values:  make(map[ValueID]*Value),
		classes: make(map[ClassID]*Class),
	}
}

func (sea *Sea) value(op Op, typ Type, imm any, args ...*Class) *Value {
	// Patch args so that if any have been replaced they point to the newest
	// thing.
	for i, a := range args {
		args[i] = a.Newest()
	}

	k := mkseakey(op, typ, imm, args...)

	v, ok := sea.m[k]
	if !ok {
		if validate := opValidation[op]; validate != nil {
			validate(typ, imm, args...)
		}

		v = sea.newValue(op, typ, imm, args...)

		// TODO: move inside newValue?
		sea.m[k] = v
	}
	return v
}

func commonClass(values iter.Seq[*Value]) *Class {
	var common *Class
	for v := range values {
		if v.class == nil {
			return nil
		}

		// TODO: this could maybe be changed to only use Newest if the first
		// common != c check fails, although it's probably not worth it.
		c := v.class.Newest()
		if common == nil {
			common = c
		}
		if common != c {
			return nil
		}
	}
	return common
}

// BUG: right now we mutate values, but really shouldn't.
// TODO: we also should probably take a func or interface to get an iterator
// from rather than hardcode map[*Value]bool, and force filtering onto the
// caller.
func (sea *Sea) class(typ Type, values map[*Value]bool) *Class {
	for v, keep := range values {
		if !keep {
			delete(values, v)
		}
	}

	if c := commonClass(maps.Keys(values)); c != nil {
		return c
	}

	// kinda gross
	var values2 []*Value
	var classes2 []*Class
	seenc := make(map[*Class]struct{})
	for v := range values {
		if c := v.class; c != nil {
			if _, ok := seenc[c]; !ok {
				seenc[c] = struct{}{}
				classes2 = append(classes2, c)
			}
		} else {
			values2 = append(values2, v)
		}
	}

	slices.SortFunc(values2, func(a, b *Value) int { return cmp.Compare(a.ID(), b.ID()) })
	slices.SortFunc(classes2, func(a, b *Class) int { return cmp.Compare(a.ID(), b.ID()) })

	return sea.newClass(typ, classes2, values2)
}

// TODO: force cloning of args onto the caller and just have newValue always
// take the ownership?
func (sea *Sea) newValue(op Op, typ Type, imm any, args ...*Class) *Value {
	sea.vid++
	id := ValueID(sea.vid)
	v := &Value{
		id:   id,
		op:   op,
		typ:  typ,
		imm:  imm,
		args: slices.Clone(args),
	}
	for _, a := range args {
		a.users[v] = struct{}{}
	}
	// sea.values[id] = v
	return v
}

func (sea *Sea) newClass(typ Type, classes []*Class, values []*Value) *Class {
	if len(classes) < 2 && len(values) == 0 {
		panic("useless")
	}

	sea.cid++
	id := ClassID(sea.cid)
	c := &Class{
		id:      id,
		typ:     typ,
		classes: classes,
		values:  values,
		users:   make(map[*Value]struct{}),
	}
	for _, o := range classes {
		o.merged = c
	}
	for _, v := range values {
		v.class = c
	}
	// sea.classes[id] = c
	return c
}
