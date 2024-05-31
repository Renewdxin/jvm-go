package rtda

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

func (self *Thread) PushFrame(frame *Frame) {
	self.stack.push(frame)
}
func (self *Thread) PopFrame() *Frame {
	return self.stack.pop()
}

func (self *Thread) CurrentFrame() *Frame {
	return self.stack.top()
}