package compiler

type Op struct{ id int32 }

// TODO: introduce OpMap[T] for efficiently mapping Op->T
var opNames = make(map[Op]string)

type ValidationFunc func(typ Type, imm any, args ...*Class)

var opValidation = make(map[Op]ValidationFunc)

var nextId int32 = 1

func DefOp(name string, validate ValidationFunc) Op {
	id := nextId
	nextId++

	// id must never be 0

	op := Op{id}
	opNames[op] = name
	opValidation[op] = validate

	return op
}

func (op Op) String() string { return opNames[op] }

// TODO: for perf we could also mark ops as commutative. That would save us from
// exploding e.g. IAdd into two instructions every time.
