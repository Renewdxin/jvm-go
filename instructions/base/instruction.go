package base

import "jvm-go/rtda"

type Instructions interface {
	FetchOperands(reader *BytecodeReader)
	Execute(frame *rtda.Frame)
}

type BranchInstruciton struct {
	Offset int
}

type Index8Instruction struct {
	Index uint
}

type Index16Indtruction struct {
	Index uint
}

type NoOperandsInstruction struct {}

func (self *BranchInstruciton) FetchOperands(reader *BytecodeReader) {
	self.Offset = int(reader.ReadInt16())
}

func (self *NoOperandsInstruction) FetchOperands(reader *BytecodeReader) {

}


func (self *Index8Instruction) FetchOperands(reader *BytecodeReader) {
	self.Index = int(reader.ReadUint8())
}

func (self *Index16Indtruction) FetchOperands(reader *BytecodeReader) {
	self.Index = int(reader.ReadUint16())
}