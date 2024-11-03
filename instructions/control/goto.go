package control

import (
	"jvm-go/instructions/base"
	"jvm-go/rtda"
)

// Branch always
type GOTO struct{ base.BranchInstruction }

func (self *GOTO) Execute(frame *rtda.Frame) {
	base.Branch(frame, self.Offset)
}
