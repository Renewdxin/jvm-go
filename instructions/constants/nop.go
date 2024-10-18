package constants

import (
	"jvm-go/instructions/base"
	"jvm-go/rtda"
)

// Do nothing
type NOP struct{ 
	base.NoOperandsInstruction 
}

func (self *NOP) Execute(frame *rtda.Frame) {
// 什么也不用做
}