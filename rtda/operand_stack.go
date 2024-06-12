package rtda

import "math"
import "jvm-go/rtda/heap"

type OperandStack struct {
	size  uint
	slots []Slot
}

func newOperandStack(maxStack uint) *OperandStack {
	if maxStack > 0 {
		return &OperandStack{
			slots: make([]Slot, maxStack),
		}
	}
	return nil
}

func (osa *OperandStack) PushInt(val int32) {
	osa.slots[osa.size].num = val
	osa.size++
}
func (osa *OperandStack) PopInt() int32 {
	osa.size--
	return osa.slots[osa.size].num
}

func (osa *OperandStack) PushFloat(val float32) {
	bits := math.Float32bits(val)
	osa.slots[osa.size].num = int32(bits)
	osa.size++
}
func (osa *OperandStack) PopFloat() float32 {
	osa.size--
	bits := uint32(osa.slots[osa.size].num)
	return math.Float32frombits(bits)
}

// long consumes two slots
func (osa *OperandStack) PushLong(val int64) {
	osa.slots[osa.size].num = int32(val)
	osa.slots[osa.size+1].num = int32(val >> 32)
	osa.size += 2
}
func (osa *OperandStack) PopLong() int64 {
	osa.size -= 2
	low := uint32(osa.slots[osa.size].num)
	high := uint32(osa.slots[osa.size+1].num)
	return int64(high)<<32 | int64(low)
}

// double consumes two slots
func (osa *OperandStack) PushDouble(val float64) {
	bits := math.Float64bits(val)
	osa.PushLong(int64(bits))
}
func (osa *OperandStack) PopDouble() float64 {
	bits := uint64(osa.PopLong())
	return math.Float64frombits(bits)
}

func (osa *OperandStack) PushRef(ref *heap.Object) {
	osa.slots[osa.size].ref = ref
	osa.size++
}
func (osa *OperandStack) PopRef() *heap.Object {
	osa.size--
	ref := osa.slots[osa.size].ref
	osa.slots[osa.size].ref = nil
	return ref
}

func (osa *OperandStack) PushSlot(slot Slot) {
	osa.slots[osa.size] = slot
	osa.size++
}
func (osa *OperandStack) PopSlot() Slot {
	osa.size--
	return osa.slots[osa.size]
}
func (osa *OperandStack) Clear() {
	osa.size = 0
	for i := range osa.slots {
		osa.slots[i].ref = nil
	}
}

func (osa *OperandStack) GetRefFromTop(n uint) *heap.Object {
	return osa.slots[osa.size-1-n].ref
}

func (osa *OperandStack) PushBoolean(val bool) {
	if val {
		osa.PushInt(1)
	} else {
		osa.PushInt(0)
	}
}
func (osa *OperandStack) PopBoolean() bool {
	return osa.PopInt() == 1
}

// todo
func NewOperandStack(maxStack uint) *OperandStack {
	return newOperandStack(maxStack)
}
