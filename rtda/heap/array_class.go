package heap

func (cl *Class) IsArray() bool {
	return cl.name[0] == '['
}

func (cl *Class) ComponentClass() *Class {
	componentClassName := getComponentClassName(cl.name)
	return cl.loader.LoadClass(componentClassName)
}

func (cl *Class) NewArray(count uint) *Object {
	if !cl.IsArray() {
		panic("Not array class: " + cl.name)
	}
	switch cl.Name() {
	case "[Z":
		return &Object{cl, make([]int8, count), nil}
	case "[B":
		return &Object{cl, make([]int8, count), nil}
	case "[C":
		return &Object{cl, make([]uint16, count), nil}
	case "[S":
		return &Object{cl, make([]int16, count), nil}
	case "[I":
		return &Object{cl, make([]int32, count), nil}
	case "[J":
		return &Object{cl, make([]int64, count), nil}
	case "[F":
		return &Object{cl, make([]float32, count), nil}
	case "[D":
		return &Object{cl, make([]float64, count), nil}
	default:
		return &Object{cl, make([]*Object, count), nil}
	}
}

func NewByteArray(loader *ClassLoader, bytes []int8) *Object {
	return &Object{loader.LoadClass("[B"), bytes, nil}
}
