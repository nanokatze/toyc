package compiler

import (
	"testing"
)

var opX = DefOp("X", nil)

func BenchmarkCreationOfValues(b *testing.B) {
	sea := NewSea()

	b.ReportAllocs()

	i := int64(0)
	for b.Loop() {
		sea.Value(OpConst, Bits32, i)
		i++
	}
}

// This actually scares me lmao
func TestGraphsWithCycles(t *testing.T) {
	sea := NewSea()

	b := Builder{Sea: sea}

	_0 := b.Value2(OpConst, Bits32, int64(0))
	x_0_0 := b.Value2(opX, Bits32, nil, _0, _0)

	sea.EquateValues(map[*Value]struct{}{_0.Value(): {}, x_0_0.Value(): {}})

	hmm := _0
	for hmm.mergedInto != nil {
		hmm = hmm.mergedInto
	}
	Dump(sea, hmm, nil)
}
