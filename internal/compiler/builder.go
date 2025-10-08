package compiler

type Builder struct {
	Sea          *Sea
	RewriteRules []RewriteRule
}

// TODO: rename to Build
func (b *Builder) Value2(op Op, typ Type, imm any, args ...*Class) *Class {
	v0 := b.Sea.Value(op, typ, imm, args...)

	values := map[*Value]struct{}{v0: {}}
	stack := []*Value{v0}
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		for _, rule := range b.RewriteRules {
			if rule.Pattern.Match(v) {
				w := rule.Replace(b.Sea, v)
				if w == nil {
					continue // rewrite didn't apply
				}
				if _, ok := values[w]; ok {
					continue // we already saw this value
				}
				values[w] = struct{}{}
				stack = append(stack, w)
			}
		}
	}

	return b.Sea.EquateValues(values)
}
