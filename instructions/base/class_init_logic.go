package base

import (
	"jvm-go/rtda"
	"jvm-go/rtda/heap"
)

// jvms 5.5

// InitClass 初始化一个类。这涉及将类标记为初始化进行中，安排其 <clinit> 方法执行，
// 以及在必要时递归初始化其超类。
func InitClass(thread *rtda.Thread, class *heap.Class) {
	// 将类标记为初始化进行中，以防止循环初始化。
	class.StartInit()

	// 安排类初始化方法 (<clinit>) 执行。
	scheduleClinit(thread, class)

	// 递归初始化超类（如果存在并且尚未初始化）。
	initSuperClass(thread, class)
}

// scheduleClinit 安排给定类的 <clinit> 方法执行。
// <clinit> 方法负责初始化静态字段和执行静态初始化块。
// 保证每个类只执行一次。
func scheduleClinit(thread *rtda.Thread, class *heap.Class) {
	// 获取类的 <clinit> 方法。
	clinit := class.GetClinitMethod()

	// 确保 clinit 方法存在并且属于此类（而不是超类）。
	// 此检查可防止循环继承情况下的无限递归。
	if clinit != nil && clinit.Class() == class {
		// 为 <clinit> 方法创建一个新的帧。
		newFrame := thread.NewFrame(clinit)

		// 将新帧推送到线程的操作数栈上。这有效地安排了 <clinit> 方法由虚拟机执行。
		thread.PushFrame(newFrame)
	}
}

// initSuperClass 初始化给定类的超类（如果存在且尚未开始初始化）。
// 这确保在子类之前初始化超类。
func initSuperClass(thread *rtda.Thread, class *heap.Class) {
	// 仅初始化非接口类的超类。接口没有传统意义上的超类（它们隐式继承自 java.lang.Object，
	// 但该初始化在别处处理）。
	if !class.IsInterface() {
		// 获取类的超类。
		superClass := class.SuperClass()

		// 如果超类存在并且尚未开始初始化，则递归初始化它。
		if superClass != nil && !superClass.InitStarted() {
			InitClass(thread, superClass)
		}
	}
}
