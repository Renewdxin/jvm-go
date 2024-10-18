package stack

import (
	"jvm-go/instructions/base"
	"jvm-go/rtda"
)

type SWAP struct {
    base.NoOperandsInstruction
}

// 交换栈顶元素
func (self *SWAP) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	slot1 := stack.PopSlot()
	slot2 := stack.PopSlot()
	stack.PushSlot(slot1)
	stack.PushSlot(slot2)
}

