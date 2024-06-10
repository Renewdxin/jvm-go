package heap

func (ob *Object) Clone() *Object {
	return &Object{
		class: ob.class,
		data:  ob.cloneData(),
	}
}

func (ob *Object) cloneData() interface{} {
	switch ob.data.(type) {
	case []int8:
		elements := ob.data.([]int8)
		elements2 := make([]int8, len(elements))
		copy(elements2, elements)
		return elements2
	case []int16:
		elements := ob.data.([]int16)
		elements2 := make([]int16, len(elements))
		copy(elements2, elements)
		return elements2
	case []uint16:
		elements := ob.data.([]uint16)
		elements2 := make([]uint16, len(elements))
		copy(elements2, elements)
		return elements2
	case []int32:
		elements := ob.data.([]int32)
		elements2 := make([]int32, len(elements))
		copy(elements2, elements)
		return elements2
	case []int64:
		elements := ob.data.([]int64)
		elements2 := make([]int64, len(elements))
		copy(elements2, elements)
		return elements2
	case []float32:
		elements := ob.data.([]float32)
		elements2 := make([]float32, len(elements))
		copy(elements2, elements)
		return elements2
	case []float64:
		elements := ob.data.([]float64)
		elements2 := make([]float64, len(elements))
		copy(elements2, elements)
		return elements2
	case []*Object:
		elements := ob.data.([]*Object)
		elements2 := make([]*Object, len(elements))
		copy(elements2, elements)
		return elements2
	default: // []Slot
		slots := ob.data.(Slots)
		slots2 := newSlots(uint(len(slots)))
		copy(slots2, slots)
		return slots2
	}
}
