package compiler

// TODO: move pattern matching machinery out of rewrite rules, it'd be useful in
// complicated passes as well, where we might need to match something but not
// (immediately) rewrite

type RewriteRule struct {
	Name    string
	Pattern *Pattern
	Binds   int
	// TODO: still can't decide if I should pass Builder or *Builder @_@
	Replace func(Builder, ...*Value) *Value
}

// TODO: make binds (*EqValues, *Value) tuples or something?
// TODO: maybe hide this and just do rule.Rewrite(v) (u, ok)? Unless we
// introduce changes in the builder that makes this impossible.
// TODO: this should actually poke a function whenever we find a match, with
// binds, rather than simply return binds. We should not just stop on whatever
// first match we find.
func (rule *RewriteRule) Match(v *Value) ([]*Value, bool) {
	m := ruleMatcher{
		binds: make([]*Value, rule.Binds),
	}
	ok := m.value(rule.Pattern, v)
	return m.binds, ok
}

// TODO: return iter.Seq?
func (rule *RewriteRule) Match2(v *Value, f func(...*Value)) {
	// call f for every match
}

// a thingy to match a single rule, that keeps track of binding and everything
type ruleMatcher struct {
	// TODO: I guess we could also keep all the possible substitutions by storing a tree
	binds []*Value
}

// TODO: ban empty classes?
func (m *ruleMatcher) class(pat *Pattern, class *EqValues) bool {
	// if pat.Op == Wildcard {
	// 	return true
	// }

	for v := range *class {
		if m.value(pat, v) {
			return true
		}
	}
	return false
}

func (m *ruleMatcher) value(pat *Pattern, v *Value) bool {
	if pat.Op != Op_ {
		if pat.Op != v.Op {
			return false
		}
		if !pat.ArgsDDD && len(v.Args) != len(pat.Args) ||
			pat.ArgsDDD && len(v.Args) < len(pat.Args) {
			return false
		}
		for i, aPat := range pat.Args {
			if !m.class(aPat, v.Args[i]) {
				return false
			}
		}
	}

	if pat.Bind != -1 {
		m.binds[pat.Bind] = v
	}
	return true
}

// TODO: implement a rewrite rule tree for faster rewrite rule matching
