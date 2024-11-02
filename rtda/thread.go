package rtda

import "jvm-go/rtda/heap"

/*
JVM
  Thread
    pc
    Stack
      Frame
        LocalVars
        OperandStack
*/
type Thread struct {
	// 当前执行的指令地址
	pc    int // the address of the instruction currently being executed
	// 栈 
	stack *Stack
	// todo
}

func (th *Thread) NewFrame(method *heap.Method) *Frame {
	return newFrame(th, method)
}

func NewThread() *Thread {
	return &Thread{
		stack: newStack(1024),
	}
}

func (th *Thread) PC() int {
	return th.pc
}
func (th *Thread) SetPC(pc int) {
	th.pc = pc
}

func (th *Thread) PushFrame(frame *Frame) {
	th.stack.push(frame)
}
func (th *Thread) PopFrame() *Frame {
	return th.stack.pop()
}

func (th *Thread) CurrentFrame() *Frame {
	return th.stack.top()
}
func (th *Thread) TopFrame() *Frame {
	return th.stack.top()
}
func (th *Thread) GetFrames() []*Frame {
	return th.stack.getFrames()
}

func (th *Thread) IsStackEmpty() bool {
	return th.stack.isEmpty()
}
func (th *Thread) ClearStack() {
	th.stack.clear()
}

