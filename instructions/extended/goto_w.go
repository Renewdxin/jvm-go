package extended

import (
	"jvm-go/instructions/base"
	"jvm-go/rtda"
)


// Branch always (wide index)
type GOTO_W struct {
	offset int
}

func (gt *GOTO_W) FetchOperands(reader *base.BytecodeReader) {
	gt.offset = int(reader.ReadInt32())
}
func (gt *GOTO_W) Execute(frame *rtda.Frame) {
	base.Branch(frame, gt.offset)
}
