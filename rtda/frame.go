package rtda

import "jvm-go/rtda/heap"

// stack frame
type Frame struct {
	lower        *Frame // stack is implemented as linked list
	localVars    LocalVars
	operandStack *OperandStack
	// 实现指令跳转
	thread       *Thread
	method       *heap.Method
	nextPC       int // the next instruction after the call
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
