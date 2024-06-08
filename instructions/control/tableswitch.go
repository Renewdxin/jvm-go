package control

import (
	"jvm-go/rtda"
	"jvm-go/instructions/base"
)

/*
tableswitch
<0-3 byte pad>
defaultbyte1
defaultbyte2
defaultbyte3
defaultbyte4
lowbyte1
lowbyte2
lowbyte3
lowbyte4
highbyte1
highbyte2
highbyte3
highbyte4
jump offsets...

Java语言中的switch-case语句有两种实现方式：
如果case值可以编码成一个索引表，则实现成tableswitch指令；
否则实现成lookupswitch指令

*/
// Access jump table by index and jump
type TABLE_SWITCH struct {
	defaultOffset int32
	low           int32
	high          int32
	jumpOffsets   []int32
}

func (swi *TABLE_SWITCH) FetchOperands(reader *base.BytecodeReader) {
	reader.SkipPadding()
	swi.defaultOffset = reader.ReadInt32()
	swi.low = reader.ReadInt32()
	swi.high = reader.ReadInt32()
	jumpOffsetsCount := swi.high - swi.low + 1
	swi.jumpOffsets = reader.ReadInt32s(jumpOffsetsCount)
}

func (swi *TABLE_SWITCH) Execute(frame *rtda.Frame) {
	index := frame.OperandStack().PopInt()

	var offset int
	if index >= swi.low && index <= swi.high {
		offset = int(swi.jumpOffsets[index-swi.low])
	} else {
		offset = int(swi.defaultOffset)
	}

	base.Branch(frame, offset)
}
