package comparisons

/**
	if<cond>指令把操作数栈顶的int变量弹出，然后跟0进行比较，满足条件则跳转。
**/

import (
	"jvm-go/rtda"
	"jvm-go/instructions/base"
)

// Branch if int comparison with zero succeeds
type IFEQ struct{ base.BranchInstruction }

func (ifeq *IFEQ) Execute(frame *rtda.Frame) {
	val := frame.OperandStack().PopInt()
	if val == 0 {
		base.Branch(frame, ifeq.Offset)
	}
}

type IFNE struct{ base.BranchInstruction }

func (ifne *IFNE) Execute(frame *rtda.Frame) {
	val := frame.OperandStack().PopInt()
	if val != 0 {
		base.Branch(frame, ifne.Offset)
	}
}

type IFLT struct{ base.BranchInstruction }

func (iflt *IFLT) Execute(frame *rtda.Frame) {
	val := frame.OperandStack().PopInt()
	if val < 0 {
		base.Branch(frame, iflt.Offset)
	}
}

type IFLE struct{ base.BranchInstruction }

func (ifle *IFLE) Execute(frame *rtda.Frame) {
	val := frame.OperandStack().PopInt()
	if val <= 0 {
		base.Branch(frame, ifle.Offset)
	}
}

type IFGT struct{ base.BranchInstruction }

func (ifgt *IFGT) Execute(frame *rtda.Frame) {
	val := frame.OperandStack().PopInt()
	if val > 0 {
		base.Branch(frame, ifgt.Offset)
	}
}

type IFGE struct{ base.BranchInstruction }

func (ifge *IFGE) Execute(frame *rtda.Frame) {
	val := frame.OperandStack().PopInt()
	if val >= 0 {
		base.Branch(frame, ifge.Offset)
	}
}
