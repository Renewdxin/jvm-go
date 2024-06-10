package extended

import (
	"jvm-go/instructions/base"
	"jvm-go/instructions/loads"
	"jvm-go/instructions/math"
	"jvm-go/instructions/stores"
	"jvm-go/rtda"
)

// 在运行时，WIDE指令会改变其后指令的行为，使其能够处理大于255的局部变量索引。

// Extend local variable index by additional bytes
type WIDE struct {
	// 被更改的指令
	modifiedInstruction base.Instruction
}

func (wi *WIDE) FetchOperands(reader *base.BytecodeReader) {
	opcode := reader.ReadUint8()
	switch opcode {
	case 0x15:
		inst := &loads.ILOAD{}
		inst.Index = uint(reader.ReadUint16())
		wi.modifiedInstruction = inst
	case 0x16:
		inst := &loads.LLOAD{}
		inst.Index = uint(reader.ReadUint16())
		wi.modifiedInstruction = inst
	case 0x17:
		inst := &loads.FLOAD{}
		inst.Index = uint(reader.ReadUint16())
		wi.modifiedInstruction = inst
	case 0x18:
		inst := &loads.DLOAD{}
		inst.Index = uint(reader.ReadUint16())
		wi.modifiedInstruction = inst
	case 0x19:
		inst := &loads.ALOAD{}
		inst.Index = uint(reader.ReadUint16())
		wi.modifiedInstruction = inst
	case 0x36:
		inst := &stores.ISTORE{}
		inst.Index = uint(reader.ReadUint16())
		wi.modifiedInstruction = inst
	case 0x37:
		inst := &stores.LSTORE{}
		inst.Index = uint(reader.ReadUint16())
		wi.modifiedInstruction = inst
	case 0x38:
		inst := &stores.FSTORE{}
		inst.Index = uint(reader.ReadUint16())
		wi.modifiedInstruction = inst
	case 0x39:
		inst := &stores.DSTORE{}
		inst.Index = uint(reader.ReadUint16())
		wi.modifiedInstruction = inst
	case 0x3a:
		inst := &stores.ASTORE{}
		inst.Index = uint(reader.ReadUint16())
		wi.modifiedInstruction = inst
	case 0x84:
		inst := &math.IINC{}
		inst.Index = uint(reader.ReadUint16())
		inst.Const = int32(reader.ReadInt16())
		wi.modifiedInstruction = inst
	case 0xa9: // ret
		panic("Unsupported opcode: 0xa9!")
	}
}

// 不改变指令操作
func (wi *WIDE) Execute(frame *rtda.Frame) {
	wi.modifiedInstruction.Execute(frame)
}
