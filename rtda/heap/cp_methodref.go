package heap

import "jvm-go/classfile"

// MethodRef 表示对方法的符号引用。它包含了对方法所属类和方法本身的引用。
type MethodRef struct {
	MemberRef         // 继承自 MemberRef，包含类名和方法描述符等信息
	method    *Method // 指向解析后的实际方法
}

// newMethodRef 创建一个新的 MethodRef 实例。
//
// 参数:
//
//	cp: 常量池指针，用于解析类名和描述符。
//	refInfo: class 文件中的常量池方法引用信息。
//
// 返回值:
//
//	*MethodRef: 新创建的 MethodRef 指针。
func newMethodRef(cp *ConstantPool, refInfo *classfile.ConstantMethodrefInfo) *MethodRef {
	ref := &MethodRef{}
	ref.cp = cp
	ref.copyMemberRefInfo(&refInfo.ConstantMemberrefInfo) // 拷贝类名和描述符等信息
	return ref
}

// ResolvedMethod 返回解析后的方法。
// 如果方法尚未解析，则调用 resolveMethodRef() 进行解析。
//
// 返回值:
//
//	*Method: 解析后的方法指针。
func (mref *MethodRef) ResolvedMethod() *Method {
	if mref.method == nil {
		mref.resolveMethodRef()
	}
	return mref.method
}

// resolveMethodRef 解析方法引用，找到对应的 Method 结构体。
// 遵循 JVM 规范 5.4.3.3。
//
// 该方法会进行以下操作：
// 1. 获取当前类 d 和被引用方法所属的类 c。
// 2. 检查 c 是否为接口，如果是则抛出 IncompatibleClassChangeError 异常。
// 3. 在 c 及其父类中查找名为 name、描述符为 descriptor 的方法。
// 4. 如果找不到方法，则抛出 NoSuchMethodError 异常。
// 5. 检查方法是否对 d 可见，如果不可见则抛出 IllegalAccessError 异常。
// 6. 将解析后的方法赋值给 mref.method。
func (mref *MethodRef) resolveMethodRef() {
	d := mref.cp.class        // 当前类
	c := mref.ResolvedClass() // 被引用方法所属的类
	if c.IsInterface() {
		panic("java.lang.IncompatibleClassChangeError") // 方法引用不能指向接口
	}

	method := lookupMethod(c, mref.name, mref.descriptor) // 在类 c 中查找方法
	if method == nil {
		panic("java.lang.NoSuchMethodError") // 找不到方法
	}
	if !method.isAccessibleTo(d) {
		panic("java.lang.IllegalAccessError") // 方法不可访问
	}

	mref.method = method // 保存解析后的方法
}

// lookupMethod 在指定的类及其父类和接口中查找方法。
//
// 参数:
//
//	class: 要查找的类。
//	name: 方法名。
//	descriptor: 方法描述符。
//
// 返回值:
//
//	*Method: 找到的方法指针，如果找不到则返回 nil。
func lookupMethod(class *Class, name, descriptor string) *Method {
	method := LookupMethodInClass(class, name, descriptor) // 先在类及其父类中查找
	if method == nil {
		method = lookupMethodInInterfaces(class.interfaces, name, descriptor) // 如果找不到，则在接口中查找
	}
	return method
}
