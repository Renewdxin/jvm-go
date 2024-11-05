package base

import (
	"fmt"
	"jvm-go/rtda"
	"jvm-go/rtda/heap"
	"strings"
)

// InvokeMethod  调用方法
// invokerFrame: 调用者的栈帧
// method:       要调用的方法
func InvokeMethod(invokerFrame *rtda.Frame, method *heap.Method) {
	// _logInvoke(callerFrame.Thread().StackDepth(), method)  // 用于调试，打印调用信息

	thread := invokerFrame.Thread()     // 获取当前线程
	newFrame := thread.NewFrame(method) // 为被调用的方法创建一个新的栈帧
	thread.PushFrame(newFrame)          // 将新的栈帧压入线程的栈

	// 传递参数
	argSlotCount := int(method.ArgSlotCount()) // 获取方法的参数槽数量
	if argSlotCount > 0 {
		// 将参数从调用者的操作数栈复制到被调用方法的局部变量表
		// 注意：参数是从后往前复制的
		for i := argSlotCount - 1; i >= 0; i-- {
			slot := invokerFrame.OperandStack().PopSlot() // 从调用者的操作数栈弹出参数
			newFrame.LocalVars().SetSlot(uint(i), slot)   // 将参数设置到被调用方法的局部变量表
		}
	}
}

// _logInvoke  打印方法调用信息 (用于调试)
func _logInvoke(stackSize uint, method *heap.Method) {
	space := strings.Repeat(" ", int(stackSize))                     // 根据栈深度生成缩进空格
	className := method.Class().Name()                               // 获取类名
	methodName := method.Name()                                      // 获取方法名
	fmt.Printf("[method]%v %v.%v()\n", space, className, methodName) // 打印调用信息
}
