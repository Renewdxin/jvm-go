package heap

import "jvm-go/classfile"


// FieldRef结构体继承了MemberRef结构体，并包含一个指向Field类型的字段。
type FieldRef struct {
	MemberRef
	field *Field
}

func newFieldRef(cp *ConstantPool, refInfo *classfile.ConstantFieldrefInfo) *FieldRef {
	ref := &FieldRef{}
	ref.cp = cp
	ref.copyMemberRefInfo(&refInfo.ConstantMemberrefInfo)
	return ref
}

// 获取字段引用所指向的字段对象
func (fr *FieldRef) ResolvedField() *Field {
	if fr.field == nil {
		fr.resolveFieldRef()
	}
	return fr.field
}

// jvms 5.4.3.2
func (fr *FieldRef) resolveFieldRef() {
	d := fr.cp.class
	c := fr.ResolvedClass()
	field := lookupField(c, fr.name, fr.descriptor)

	if field == nil {
		panic("java.lang.NoSuchFieldError")
	}
	if !field.isAccessibleTo(d) {
		panic("java.lang.IllegalAccessError")
	}

	fr.field = field
}

// 在给定的类中查找与指定名称和描述符匹配的字段对象
func lookupField(c *Class, name, descriptor string) *Field {
	for _, field := range c.fields {
		if field.name == name && field.descriptor == descriptor {
			return field
		}
	}

	for _, iface := range c.interfaces {
		if field := lookupField(iface, name, descriptor); field != nil {
			return field
		}
	}

	if c.superClass != nil {
		return lookupField(c.superClass, name, descriptor)
	}

	return nil
}
