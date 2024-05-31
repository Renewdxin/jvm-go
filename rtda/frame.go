package rtda

type Frame struct {
	lower *Frame
	localVars localVars
	operandStack *OperandStack
}

func newFrame(maxLocals, maxStack uint) *Frame {
	return &Frame{
        localVars: newLocalVars(maxLocals),
        operandStack: newOperandStack(maxStack),
    }
}

