package extended

import (
	"jvm-go/instructions/base"
	"jvm-go/rtda"
)


// Branch always (wide index)
type GOTO_W struct {
	offset int
}

// goto_w指令和goto指令的唯一区别就是索引从2字节变成了4字节
func (gt *GOTO_W) FetchOperands(reader *base.BytecodeReader) {
	gt.offset = int(reader.ReadInt32())
}
func (gt *GOTO_W) Execute(frame *rtda.Frame) {
	base.Branch(frame, gt.offset)
}
