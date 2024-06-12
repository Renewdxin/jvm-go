package heap

// symbolic reference 类符号引用
type SymRef struct {
	cp        *ConstantPool
	className string
	class     *Class
}

func (sr *SymRef) ResolvedClass() *Class {
	if sr.class == nil {
		sr.resolveClassRef()
	}
	return sr.class
}

// jvms8 5.4.3.1
func (sr *SymRef) resolveClassRef() {
	d := sr.cp.class
	c := d.loader.LoadClass(sr.className)
	if !c.isAccessibleTo(d) {
		panic("java.lang.IllegalAccessError")
	}

	sr.class = c
}
