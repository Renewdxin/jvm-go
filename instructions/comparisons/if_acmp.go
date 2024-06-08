package comparisons

import (
	"jvm-go/rtda"
	"jvm-go/instructions/base"
)

// Branch if reference comparison succeeds
type IF_ACMPEQ struct{ base.BranchInstruction }

func (quote *IF_ACMPEQ) Execute(frame *rtda.Frame) {
	if _acmp(frame) {
		base.Branch(frame, quote.Offset)
	}
}

type IF_ACMPNE struct{ base.BranchInstruction }

func (quote *IF_ACMPNE) Execute(frame *rtda.Frame) {
	if !_acmp(frame) {
		base.Branch(frame, quote.Offset)
	}
}

func _acmp(frame *rtda.Frame) bool {
	stack := frame.OperandStack()
	ref2 := stack.PopRef()
	ref1 := stack.PopRef()
	return ref1 == ref2 // todo
}
