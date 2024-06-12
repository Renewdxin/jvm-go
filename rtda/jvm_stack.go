package rtda

// jvm stack
type Stack struct {
	maxSize uint
	size    uint
	_top    *Frame // stack is implemented as linked list
}

func newStack(maxSize uint) *Stack {
	return &Stack{
		maxSize: maxSize,
	}
}

func (sta *Stack) push(frame *Frame) {
	if sta.size >= sta.maxSize {
		panic("java.lang.StackOverflowError")
	}

	if sta._top != nil {
		frame.lower = sta._top
	}

	sta._top = frame
	sta.size++
}

func (sta *Stack) pop() *Frame {
	if sta._top == nil {
		panic("jvm stack is empty!")
	}

	top := sta._top
	sta._top = top.lower
	top.lower = nil
	sta.size--

	return top
}

func (sta *Stack) top() *Frame {
	if sta._top == nil {
		panic("jvm stack is empty!")
	}

	return sta._top
}

func (sta *Stack) getFrames() []*Frame {
	frames := make([]*Frame, 0, sta.size)
	for frame := sta._top; frame != nil; frame = frame.lower {
		frames = append(frames, frame)
	}
	return frames
}

func (sta *Stack) isEmpty() bool {
	return sta._top == nil
}

func (sta *Stack) clear() {
	for !sta.isEmpty() {
		sta.pop()
	}
}
