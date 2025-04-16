package heap

type Object struct {
	// 类
	class *Class
	// 数据
	data interface{} // Slots for Object, []int32 for int[] ...
	// 额外数据
	extra interface{}
}

// create normal (non-array) object
func newObject(class *Class) *Object {
	return &Object{
		class: class,
		data:  newSlots(class.instanceSlotCount),
	}
}

// getters & setters
func (ob *Object) Class() *Class {
	return ob.class
}
func (ob *Object) Data() interface{} {
	return ob.data
}
func (ob *Object) Fields() Slots {
	return ob.data.(Slots)
}
func (ob *Object) Extra() interface{} {
	return ob.extra
}
func (ob *Object) SetExtra(extra interface{}) {
	ob.extra = extra
}

func (ob *Object) IsInstanceOf(class *Class) bool {
	return class.IsAssignableFrom(ob.class)
}

// reflection
func (ob *Object) GetRefVar(name, descriptor string) *Object {
	field := ob.class.getField(name, descriptor, false)
	slots := ob.data.(Slots)
	return slots.GetRef(field.slotId)
}
func (ob *Object) SetRefVar(name, descriptor string, ref *Object) {
	field := ob.class.getField(name, descriptor, false)
	slots := ob.data.(Slots)
	slots.SetRef(field.slotId, ref)
}
func (ob *Object) SetIntVar(name, descriptor string, val int32) {
	field := ob.class.getField(name, descriptor, false)
	slots := ob.data.(Slots)
	slots.SetInt(field.slotId, val)
}
func (ob *Object) GetIntVar(name, descriptor string) int32 {
	field := ob.class.getField(name, descriptor, false)
	slots := ob.data.(Slots)
	return slots.GetInt(field.slotId)
}
