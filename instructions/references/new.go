package references

import (
	"jvm-go/instructions/base"
	"jvm-go/rtda"
	"jvm-go/rtda/heap"
)

// NEW 指令用于创建一个新的对象实例。
type NEW struct {
	base.Index16Instruction
}

// Execute 执行 NEW 指令。
// frame: 当前栈帧
func (ne *NEW) Execute(frame *rtda.Frame) {
	// 1. 获取常量池
	cp := frame.Method().Class().ConstantPool()

	// 2. 获取 ClassRef
	// 从常量池中获取索引为 ne.Index 的常量，并将其转换为 ClassRef 类型。
	// 这个 ClassRef 指向要创建实例的类。
	classRef := cp.GetConstant(ne.Index).(*heap.ClassRef)

	// 3. 解析类
	// 获取 ClassRef 指向的类。如果类还没有被加载，则会触发类加载过程。
	class := classRef.ResolvedClass()

	// 4. 检查类是否已初始化
	// 如果类还没有被初始化，则先初始化类，并暂停当前方法的执行，直到类初始化完成。
	if !class.InitStarted() {
		// 恢复下一条指令的地址，以便在类初始化完成后继续执行当前方法。
		frame.RevertNextPC()
		// 初始化类。这会执行类的 <clinit> 方法。
		base.InitClass(frame.Thread(), class)
		// 返回，暂停当前方法的执行。
		return
	}

	// 5. 检查类是否为接口或抽象类
	// 如果类是接口或抽象类，则抛出 InstantiationError 异常。
	if class.IsInterface() || class.IsAbstract() {
		panic("java.lang.InstantiationError")
	}

	// 6. 创建对象实例
	// 创建一个新的对象实例。
	ref := class.NewObject()

	// 7. 将对象实例推入操作数栈
	// 将新创建的对象实例的引用推入操作数栈。
	frame.OperandStack().PushRef(ref)
}
