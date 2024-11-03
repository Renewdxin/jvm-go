package heap

// symbolic reference 类符号引用  ——  注释说明 SymRef 结构体的作用
type SymRef struct { // 定义 SymRef 结构体
	cp        *ConstantPool // 指向常量池的指针，符号引用需要从常量池中获取信息
	className string        // 存储符号引用的类名（以完全限定名的形式，例如 "java/lang/String"）
	class     *Class        // 指向已解析的类的指针，初始值为 nil，解析后指向实际的 Class 对象
}

// ResolvedClass 方法返回已解析的 Class 对象。  ——  注释说明方法的作用
func (sr *SymRef) ResolvedClass() *Class { // SymRef 的方法，用于获取已解析的 Class 对象
	if sr.class == nil { // 检查是否已经解析过
		sr.resolveClassRef() // 如果未解析，则调用 resolveClassRef 方法进行解析
	}
	return sr.class // 返回已解析的 Class 对象
}

// jvms8 5.4.3.1  ——  注释说明此方法实现的 JVM 规范
func (sr *SymRef) resolveClassRef() { // SymRef 的方法，用于解析符号引用
	d := sr.cp.class                      // 获取当前类（定义 SymRef 的类）
	c := d.loader.LoadClass(sr.className) // 使用类加载器加载符号引用指向的类
	if !c.isAccessibleTo(d) {             // 检查已加载的类 c 是否对当前类 d 可访问
		panic("java.lang.IllegalAccessError") // 如果不可访问，则抛出 IllegalAccessError 异常
	}

	sr.class = c // 将已解析的类对象赋值给 sr.class 字段
}
