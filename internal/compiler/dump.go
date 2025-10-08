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
	fmt.Fprintf(&line, "%v %v", f(class), class.Type())
	if len(class.Values) > 1 {
		fmt.Fprintln(&line)
	}
	for v := range class.Values {
		if len(class.Values) > 1 {
			fmt.Fprint(&line, "    = ")
		} else {
			fmt.Fprint(&line, " = ")
		}
		fmt.Fprintf(&line, "%v", v.Op())
		if imm := v.Imm(); imm != nil {
			fmt.Fprintf(&line, " %v", imm)
		}
		for _, a := range v.Args {
			fmt.Fprintf(&line, " %v", f(a))
			dump(sea, a, dumped, f)
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
		f = func(c *Class) string { return c.ID.String() }
	}
	dumped := make(map[*Class]struct{})
	dump(sea, class, dumped, f)
}
