package compiler

// TODO: implement binding. That's actually not very trivial with egraphs at
// play

type Pattern struct {
	Op      Op
	Args    []*Pattern
	ArgsDDD bool
	// Bind    int
}

func (p *Pattern) MatchClass(c *Class) bool {
	for n := range c.Values {
		if p.Match(n) {
			return true
		}
	}
	return false
}

func (p *Pattern) Match(v *Value) bool {
	if p.Op != (Op{}) {
		if v.Op() != p.Op {
			return false
		}
		if !p.ArgsDDD && len(v.Args) != len(p.Args) ||
			p.ArgsDDD && len(v.Args) < len(p.Args) {
			return false
		}
		for i, r := range p.Args {
			a := v.Args[i]
			if !r.MatchClass(a) {
				return false
			}
		}
	}
	return true
}
