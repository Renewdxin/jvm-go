package extended

import (
	"jvm-go/instructions/base"
	"jvm-go/rtda"
)

// Branch if reference is null
type IFNULL struct{ base.BranchInstruction }

func (inu *IFNULL) Execute(frame *rtda.Frame) {
	ref := frame.OperandStack().PopRef()
	if ref == nil {
		base.Branch(frame, inu.Offset)
	}
}

// Branch if reference not null
type IFNONNULL struct{ base.BranchInstruction }

func (inu *IFNONNULL) Execute(frame *rtda.Frame) {
	ref := frame.OperandStack().PopRef()
	if ref != nil {
		base.Branch(frame, inu.Offset)
	}
}
