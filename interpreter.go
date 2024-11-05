package main

import (
	"fmt"
	"jvm-go/instructions"
	"jvm-go/instructions/base"
	"jvm-go/rtda"
)

// interpret 解释执行字节码。
// thread: 当前线程
// logInst: 是否打印指令执行信息
func interpret(thread *rtda.Thread, logInst bool) {
	// 使用 defer 和 recover 机制捕获运行时错误，并在发生错误时打印栈帧信息。
	defer catchErr(thread)
	loop(thread, logInst)
}

// catchErr 捕获 panic 并打印栈帧信息。
// thread: 当前线程
func catchErr(thread *rtda.Thread) {
	if r := recover(); r != nil {
		logFrames(thread)
		panic(r) // 重新抛出异常，以便上层处理
	}
}

// loop 解释执行循环。
// thread: 当前线程
// logInst: 是否打印指令执行信息
func loop(thread *rtda.Thread, logInst bool) {
	reader := &base.BytecodeReader{} // 字节码读取器
	for {
		frame := thread.CurrentFrame() // 获取当前栈帧
		pc := frame.NextPC()           // 获取下一条指令的地址
		thread.SetPC(pc)               // 设置线程的程序计数器

		// 解码指令
		reader.Reset(frame.Method().Code(), pc)     // 重置字节码读取器
		opcode := reader.ReadUint8()                // 读取操作码
		inst := instructions.NewInstruction(opcode) // 创建指令
		inst.FetchOperands(reader)                  // 读取操作数
		frame.SetNextPC(reader.PC())                // 更新下一条指令的地址

		// 打印指令执行信息（如果开启了日志）
		if logInst {
			logInstruction(frame, inst)
		}

		// 执行指令
		inst.Execute(frame)

		// 如果操作数栈为空，则退出循环（程序执行结束）
		if thread.IsStackEmpty() {
			break
		}
	}
}

// logInstruction 打印指令执行信息。
// frame: 当前栈帧
// inst: 当前指令
func logInstruction(frame *rtda.Frame, inst base.Instruction) {
	method := frame.Method()
	className := method.Class().Name()
	methodName := method.Name()
	pc := frame.Thread().PC()
	fmt.Printf("%v.%v() #%2d %T %v\n", className, methodName, pc, inst, inst)
}

// logFrames 打印栈帧信息。
// thread: 当前线程
func logFrames(thread *rtda.Thread) {
	// 遍历所有栈帧并打印信息
	for !thread.IsStackEmpty() {
		frame := thread.PopFrame()                      // 弹出栈帧
		method := frame.Method()                        // 获取方法
		className := method.Class().Name()              // 获取类名
		lineNum := method.GetLineNumber(frame.NextPC()) // 获取行号
		fmt.Printf(">> line:%4d pc:%4d %v.%v%v \n",
			lineNum, frame.NextPC(), className, method.Name(), method.Descriptor())
	}
}
