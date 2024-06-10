package heap

import "jvm-go/classfile"

type Field struct {
	ClassMember
	constValueIndex uint
	slotId          uint
}

func newFields(class *Class, cfFields []*classfile.MemberInfo) []*Field {
	fields := make([]*Field, len(cfFields))
	for i, cfField := range cfFields {
		fields[i] = &Field{}
		fields[i].class = class
		fields[i].copyMemberInfo(cfField)
		fields[i].copyAttributes(cfField)
	}
	return fields
}
func (fi *Field) copyAttributes(cfField *classfile.MemberInfo) {
	if valAttr := cfField.ConstantValueAttribute(); valAttr != nil {
		fi.constValueIndex = uint(valAttr.ConstantValueIndex())
	}
}

func (fi *Field) IsVolatile() bool {
	return 0 != fi.accessFlags&ACC_VOLATILE
}
func (fi *Field) IsTransient() bool {
	return 0 != fi.accessFlags&ACC_TRANSIENT
}
func (fi *Field) IsEnum() bool {
	return 0 != fi.accessFlags&ACC_ENUM
}

func (fi *Field) ConstValueIndex() uint {
	return fi.constValueIndex
}
func (fi *Field) SlotId() uint {
	return fi.slotId
}
func (fi *Field) isLongOrDouble() bool {
	return fi.descriptor == "J" || fi.descriptor == "D"
}

// reflection
func (fi *Field) Type() *Class {
	className := toClassName(fi.descriptor)
	return fi.class.loader.LoadClass(className)
}
