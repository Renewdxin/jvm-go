package stack

import (
	"jvm-go/instructions/base"
	"jvm-go/rtda"
)

type DUP struct {
	base.NoOperandsInstruction
}

type DUP_X1 struct {
	base.NoOperandsInstruction
}

type DUP_X2 struct {
	base.NoOperandsInstruction
}

type DUP2 struct {
	base.NoOperandsInstruction
}

type DUP2_X1 struct {
	base.NoOperandsInstruction
}

type DUP2_X2 struct {
	base.NoOperandsInstruction
}

func(self *DUP) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	slot := stack.PopSlot()
	stack.PushSlot(slot)
	stack.PushSlot(slot)
}

func(self *DUP_X1) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	slot := stack.PopSlot()
	stack.PushSlot(slot)
	
}