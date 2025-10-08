package compiler

type Op struct {
	Name        string
	Commutative bool // TODO: shove bools into flags?
	// Validate func(v *Value) error
}

func (op *Op) String() string { return op.Name }

// TODO: make a function for constructing ops so we can record where it's
// defined etc?
var (
	OpConst = &Op{Name: "Const"}

	OpTupleMake    = &Op{Name: "TupleMake"}
	OpTupleExtract = &Op{Name: "TupleExtract"}

	// TODO: move anything but the most basic ops into a separate package?
	OpIAdd = &Op{Name: "IAdd", Commutative: true}
	OpIMul = &Op{Name: "IMul", Commutative: true}
)
