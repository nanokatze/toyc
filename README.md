This is the second toy compiler in my series of toy compiler, building
upon [cg2c](https://github.com/nanokatze/cg2c).

Cg2c had a bunch of interesting things going for it, notably
* Use of RVSDG instead of CFG for control flow
* Lacking an explicit schedule of instructions, making it easy to insert
  instructions and trivially enabling dead code elimination
* Immutable instructions, which in turn enable common subexpression elimination
  and rule-based rewriting at instruction creation time, rather than needing to
  apply a pass

Compared to cg2c, the new Big Things this toy compiler has going for it are
* Ad-hoc extensibility. Dependent code can add new ops and types.
* E-graphs. The compiler maintains all the different program representations
  that appear throughout the compilation process. The e-graphs, surprisingly,
  are turning out to be a rather small change to the compiler that's otherwise
  constructed in a way very similar to cg2c.

At the moment, this compiler can not represent control flow.
