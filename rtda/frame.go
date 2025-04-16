package rtda

import "jvm-go/rtda/heap"

// stack frame
type Frame struct {
	// 用于实现链表结构
	lower *Frame
	// 局部变量表
	localVars LocalVars
	// 操作数栈
	operandStack *OperandStack
	// 线程
	thread *Thread
	// 方法
	method *heap.Method
	// 下一条指令的地址
	nextPC int
}

func newFrame(thread *Thread, method *heap.Method) *Frame {
	return &Frame{
		thread:       thread,
		method:       method,
		localVars:    newLocalVars(method.MaxLocals()),
		operandStack: newOperandStack(method.MaxStack()),
	}
}

// getters & setters
func (fra *Frame) LocalVars() LocalVars {
	return fra.localVars
}
func (fra *Frame) OperandStack() *OperandStack {
	return fra.operandStack
}
func (fra *Frame) Thread() *Thread {
	return fra.thread
}
func (fra *Frame) Method() *heap.Method {
	return fra.method
}
func (fra *Frame) NextPC() int {
	return fra.nextPC
}
func (fra *Frame) SetNextPC(nextPC int) {
	fra.nextPC = nextPC
}

func (fra *Frame) RevertNextPC() {
	fra.nextPC = fra.thread.pc
}
