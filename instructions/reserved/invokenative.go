package reserved

import (
	"jvm-go/instructions/base"
	"jvm-go/native"
	"jvm-go/rtda"
)

// INVOKE_NATIVE 指令调用本地方法
type INVOKE_NATIVE struct{ base.NoOperandsInstruction } // 没有操作数

// Execute 方法执行 INVOKE_NATIVE 指令
func (n *INVOKE_NATIVE) Execute(frame *rtda.Frame) {
	method := frame.Method()                // 获取当前正在执行的方法
	className := method.Class().Name()      // 获取当前方法所属类的名称
	methodName := method.Name()             // 获取当前方法的名称
	methodDescriptor := method.Descriptor() // 获取当前方法的描述符

	// 根据类名、方法名和方法描述符查找本地方法
	// 如果找不到，则抛出UnsatisfiedLinkError异常，否则直接调用本地方法
	nativeMethod := native.FindNativeMethod(className, methodName, methodDescriptor)
	if nativeMethod == nil { // 如果找不到本地方法，则抛出 UnsatisfiedLinkError 异常
		methodInfo := className + "." + methodName + methodDescriptor
		panic("java.lang.UnsatisfiedLinkError: " + methodInfo)
	}

	// 调用找到的本地方法
	nativeMethod(frame)
}
