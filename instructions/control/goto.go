package control

import (
	"jvm-go/rtda"
	"jvm-go/instructions/base"
)

// goto 指令 无条件跳转

// Branch always
type GOTO struct{ base.BranchInstruction }

func (g *GOTO) Execute(frame *rtda.Frame) {
	base.Branch(frame, g.Offset)
}
