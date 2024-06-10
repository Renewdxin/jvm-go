package heap

import "math"

type Slot struct {
	num int32
	ref *Object
}

type Slots []Slot

func newSlots(slotCount uint) Slots {
	if slotCount > 0 {
		return make([]Slot, slotCount)
	}
	return nil
}

func (sl Slots) SetInt(index uint, val int32) {
	sl[index].num = val
}
func (sl Slots) GetInt(index uint) int32 {
	return sl[index].num
}

func (sl Slots) SetFloat(index uint, val float32) {
	bits := math.Float32bits(val)
	sl[index].num = int32(bits)
}
func (sl Slots) GetFloat(index uint) float32 {
	bits := uint32(sl[index].num)
	return math.Float32frombits(bits)
}

// long consumes two slots
func (sl Slots) SetLong(index uint, val int64) {
	sl[index].num = int32(val)
	sl[index+1].num = int32(val >> 32)
}
func (sl Slots) GetLong(index uint) int64 {
	low := uint32(sl[index].num)
	high := uint32(sl[index+1].num)
	return int64(high)<<32 | int64(low)
}

// double consumes two slots
func (sl Slots) SetDouble(index uint, val float64) {
	bits := math.Float64bits(val)
	sl.SetLong(index, int64(bits))
}
func (sl Slots) GetDouble(index uint) float64 {
	bits := uint64(sl.GetLong(index))
	return math.Float64frombits(bits)
}

func (sl Slots) SetRef(index uint, ref *Object) {
	sl[index].ref = ref
}
func (sl Slots) GetRef(index uint) *Object {
	return sl[index].ref
}
