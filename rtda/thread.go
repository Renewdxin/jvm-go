package rtda

import "jvm-go/rtda/heap"

type Thread struct {
	pc int
	// pointer to stack
	stack *Stack
}

func NewThread() *Thread {
    return &Thread{
		stack: newStack(1024),
	}
}

func (self *Thread) NewFrame(method *heap.Method) *Frame {
	return newFrame(self, method)
}

func (self *Thread) PushFrame(frame *Frame) {
	self.stack.push(frame)
}
func (self *Thread) PopFrame() *Frame {
	return self.stack.pop()
}

func (self *Thread) CurrentFrame() *Frame {
	return self.stack.top()
}

func (self *Thread) PC() int{
	return self.pc
}