package heap

import "jvm-go/classfile"

type InterfaceMethodRef struct {
	MemberRef
	method *Method
}

func newInterfaceMethodRef(cp *ConstantPool, refInfo *classfile.ConstantInterfaceMethodrefInfo) *InterfaceMethodRef {
	ref := &InterfaceMethodRef{}
	ref.cp = cp
	ref.copyMemberRefInfo(&refInfo.ConstantMemberrefInfo)
	return ref
}

func (imref *InterfaceMethodRef) ResolvedInterfaceMethod() *Method {
	if imref.method == nil {
		imref.resolveInterfaceMethodRef()
	}
	return imref.method
}

// jvms8 5.4.3.4
func (imref *InterfaceMethodRef) resolveInterfaceMethodRef() {
	d := imref.cp.class
	c := imref.ResolvedClass()
	if !c.IsInterface() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	method := lookupInterfaceMethod(c, imref.name, imref.descriptor)
	if method == nil {
		panic("java.lang.NoSuchMethodError")
	}
	if !method.isAccessibleTo(d) {
		panic("java.lang.IllegalAccessError")
	}

	imref.method = method
}

// todo
func lookupInterfaceMethod(iface *Class, name, descriptor string) *Method {
	for _, method := range iface.methods {
		if method.name == name && method.descriptor == descriptor {
			return method
		}
	}

	return lookupMethodInInterfaces(iface.interfaces, name, descriptor)
}
