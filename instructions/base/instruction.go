package base

import "jvm-go/rtda"

type Instruction interface {
	FetchOperands(reader *BytecodeReader)
	Execute(frame *rtda.Frame)
}

type NoOperandsInstruction struct {
	// empty
}

func (npi *NoOperandsInstruction) FetchOperands(reader *BytecodeReader) {
	// nothing to do
}

type BranchInstruction struct {
	Offset int
}

func (npi *BranchInstruction) FetchOperands(reader *BytecodeReader) {
	npi.Offset = int(reader.ReadInt16())
}

type Index8Instruction struct {
	Index uint
}

func (npi *Index8Instruction) FetchOperands(reader *BytecodeReader) {
	npi.Index = uint(reader.ReadUint8())
}

type Index16Instruction struct {
	Index uint
}

func (npi *Index16Instruction) FetchOperands(reader *BytecodeReader) {
	npi.Index = uint(reader.ReadUint16())
}
