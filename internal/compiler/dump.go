package compiler

import (
	"fmt"
	"os"
	"strings"
)

// TODO: functions to dump the entire sea, dump the induced subgraph, dump a
// single class, etc.

func dump(sea *Sea, class *Class, dumped map[*Class]struct{}, f func(*Class) string) {
	if _, ok := dumped[class]; ok {
		return
	}
	dumped[class] = struct{}{}

	var line strings.Builder
	fmt.Fprintf(&line, "%v %v\n", f(class), class.Type())
	for _, c := range class.classes {
		fmt.Fprintf(&line, " = %v\n", c.id)
		dump(sea, c, dumped, f)
	}
	for _, v := range class.values {
		fmt.Fprint(&line, " = ")
		fmt.Fprintf(&line, "%v", v.Op())
		if imm := v.Imm(); imm != nil {
			fmt.Fprintf(&line, " %v", imm)
		}
		for _, a := range v.Args() {
			fmt.Fprintf(&line, " %v", f(a))
			dump(sea, a, dumped, f)
		}
		if v.class != class {
			fmt.Fprintf(&line, " (from %v)", v.class.id)
		}
		fmt.Fprintf(&line, "\n")
	}

	fmt.Fprint(os.Stderr, line.String())
}

// TODO: accept some kind of object which would control how class IDs are
// printed. This would allow to print associated data (e.g. regalloc) etc.
// TODO: allow user to omit types.
func Dump(sea *Sea, class *Class, f func(*Class) string) {
	if f == nil {
		f = func(c *Class) string { return c.id.String() }
	}
	dumped := make(map[*Class]struct{})
	dump(sea, class, dumped, f)
}
