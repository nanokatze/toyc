package compiler

// special op for use in pattern matching
var Op_ = &Op{Name: "_"}

type Pattern struct {
	Op      *Op
	Args    []*Pattern
	ArgsDDD bool
	Bind    int
}
