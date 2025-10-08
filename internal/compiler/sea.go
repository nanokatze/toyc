package compiler

import (
	"iter"
	"maps"
	"slices"
	"strconv"
)

type Type interface {
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
	Args  []*Class // TODO: hide this
	class *Class
	// deleted bool
}

func (v *Value) ID() ValueID { return v.id }

func (v *Value) Op() Op { return v.op }

func (v *Value) Type() Type { return v.typ }

func (v *Value) Imm() any { return v.imm }

// Class that this value is in. May be nil if the value is yet to be placed into
// any class, or has been deleted.
func (v *Value) Class() *Class { return v.class }

type ClassID int32

func (id ClassID) String() string { return "c" + strconv.FormatInt(int64(id), 10) }

// Class is a set of equivalent value definitions.
type Class struct {
	ID         ClassID
	typ        Type
	Values     map[*Value]struct{} // TODO: hide and replace with a method to return iter.Seq[*Value]
	users      map[*Value]struct{}
	mergedInto *Class
}

/*
func (c *Class) replacedBy() *Class {
	for c.replacedBy != nil {
		c = c.replacedBy
	}
	return c
}
*/

func (c *Class) Type() Type {
	return c.typ
}

func (to *Class) add(v *Value) {
	if v.class != nil {
		panic("trying to add a value that is already in a class")
	}
	to.Values[v] = struct{}{}
	v.class = to
}

func (to *Class) moveValues(from *Class) {
	for v := range from.Values {
		if v.class != from {
			panic("inconsistent sea")
		}
		delete(from.Values, v)
		to.Values[v] = struct{}{}
		v.class = to
	}
}

// TODO: delete this
func (c *Class) Value() *Value {
	if len(c.Values) != 1 {
		panic("class must contain exactly one value")
	}
	for v := range c.Values {
		return v
	}
	panic("unreachable")
}

func (c *Class) Values2() iter.Seq[*Value] {
	return maps.Keys(c.Values)
}

// TODO: make sure this is the interface we need actually
// TODO: do we want this to be called Uses or Users?
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

func (sea *Sea) Value(op Op, typ Type, imm any, args ...*Class) *Value {
	// TODO: patch args so that if any have been replaced they point to the
	// newest thing.

	k := mkseakey(op, typ, imm, args...)

	v, ok := sea.m[k]
	if !ok {
		if validate := opValidation[op]; validate != nil {
			validate(typ, imm, args...)
		}

		v = sea.newValue(op, typ, imm, args...)

		sea.postArgsChange(v)
	}
	return v
}

// TODO: this could take an iter.Seq[*Value] ig
func assertValuesHaveTheSameTypes(values map[*Value]struct{}) Type {
	var typ Type
	for v := range values {
		t := v.Type()
		if typ != nil && typ != t {
			panic("wat")
		}
		typ = t
	}
	if typ == nil {
		panic("guh")
	}
	return typ
}

// TODO: make this accept a func that produces an iter I guess?
func (sea *Sea) EquateValues(values map[*Value]struct{}) *Class {
	var c *Class
	for v := range values {
		// Pick the class with the most values in it so we hopefully have less
		// work to do to merge.
		if v.class != nil && (c == nil || len(c.Values) < len(v.class.Values)) {
			c = v.class
		}
	}
	if c == nil {
		typ := assertValuesHaveTheSameTypes(values)
		c = sea.newClass(typ)
	}

	if c.mergedInto != nil {
		panic("unreachable") // TODO: explain
	}

	for v := range values {
		sea.EquateValue(c, v)
	}

	return c
}

// Equate v to all values in c.
//
// TODO: I guess we could make this public as well?
func (sea *Sea) EquateValue(c *Class, v *Value) {
	if _, ok := c.Values[v]; ok {
		if v.class != c {
			panic("inconsistent sea")
		}
		return // already in the class
	}

	if v.class != nil {
		sea.equateClasses(c, v.class)
	} else {
		c.add(v)
	}
}

// Equate values in c1 and c2 by merging c2 into c1. Values from c2 are moved to
// c1 and all uses of c2 are replaced with uses of c1. Values that use to c2 may
// be replaced by already existing values that are equivalent but use c1 in
// place of c2.
//
// TODO: ponder making this public
func (sea *Sea) equateClasses(c1, c2 *Class) {
	if c1 == c2 {
		return
	}

	if c1.Type() != c2.Type() {
		panic("types differ")
	}

	if c1.mergedInto != nil {
		// TODO: actually just chase pointers at the start of the func instead?
		panic("can't merge into a class that has been merged")
	}

	if c2.mergedInto != nil {
		return
	}
	c2.mergedInto = c1

	c1.moveValues(c2)

	// Replace uses of c2 with uses of c1.
	for u := range c2.Users() {
		sea.preArgsChange(u)

		for i := range u.Args {
			if u.Args[i] == c2 {
				u.Args[i] = c1
			}
		}

		sea.postArgsChange(u)
	}
}

// TODO: rename this and the next one to something more sensible
func (sea *Sea) preArgsChange(v *Value) {
	for _, a := range v.Args {
		delete(a.users, v)
	}

	key := mkseakey(v.op, v.typ, v.imm, v.Args...)
	delete(sea.m, key)
}

func (sea *Sea) postArgsChange(v *Value) {
	key := mkseakey(v.op, v.typ, v.imm, v.Args...)

	if v0, ok := sea.m[key]; ok {
		// There's an equivalent value already, delete v and replace it with v0
		// in its class.

		if c := v.class; c != nil {
			delete(c.Values, v)
			// TODO: introduce deletion marker in value?

			sea.EquateValue(c, v0)
		}
	} else {
		sea.m[key] = v

		for _, a := range v.Args {
			a.users[v] = struct{}{}
		}
	}
}

// TODO: probably delete this garbage
func (sea *Sea) KillValue(v *Value) {
	for _, a := range v.Args {
		delete(a.users, v)
	}

	if c := v.class; c != nil {
		delete(c.Values, v)
	}
}

// TODO: add a method for removing values from a class? Would be useful when
// trimming the Sea.

func (sea *Sea) newValue(op Op, typ Type, imm any, args ...*Class) *Value {
	sea.vid++
	id := ValueID(sea.vid)
	v := &Value{
		id:   id,
		op:   op,
		typ:  typ,
		imm:  imm,
		Args: slices.Clone(args),
	}
	// sea.values[id] = v
	return v
}

func (sea *Sea) newClass(typ Type) *Class {
	sea.cid++
	id := ClassID(sea.cid)
	c := &Class{
		ID:     id,
		typ:    typ,
		Values: make(map[*Value]struct{}),
		users:  make(map[*Value]struct{}),
	}
	// sea.classes[id] = c
	return c
}
