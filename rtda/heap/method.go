package heap

import "jvm-go/classfile"

type Method struct {
	ClassMember
	// 存放操作数栈
	maxStack                uint
	// 存放局部变量表
	maxLocals               uint
	// 存放方法字节码
	code                    []byte
	exceptionTable          ExceptionTable // todo: rename
	lineNumberTable         *classfile.LineNumberTableAttribute
	exceptions              *classfile.ExceptionsAttribute // todo: rename
	parameterAnnotationData []byte                         // RuntimeVisibleParameterAnnotations_attribute
	annotationDefaultData   []byte                         // AnnotationDefault_attribute
	parsedDescriptor        *MethodDescriptor
	argSlotCount            uint
}

func newMethods(class *Class, cfMethods []*classfile.MemberInfo) []*Method {
	methods := make([]*Method, len(cfMethods))
	for i, cfMethod := range cfMethods {
		methods[i] = newMethod(class, cfMethod)
	}
	return methods
}

func newMethod(class *Class, cfMethod *classfile.MemberInfo) *Method {
	method := &Method{}
	method.class = class
	method.copyMemberInfo(cfMethod)
	method.copyAttributes(cfMethod)
	md := parseMethodDescriptor(method.descriptor)
	method.parsedDescriptor = md
	method.calcArgSlotCount(md.parameterTypes)
	if method.IsNative() {
		method.injectCodeAttribute(md.returnType)
	}
	return method
}

func (me *Method) copyAttributes(cfMethod *classfile.MemberInfo) {
	if codeAttr := cfMethod.CodeAttribute(); codeAttr != nil {
		me.maxStack = codeAttr.MaxStack()
		me.maxLocals = codeAttr.MaxLocals()
		me.code = codeAttr.Code()
		me.lineNumberTable = codeAttr.LineNumberTableAttribute()
		me.exceptionTable = newExceptionTable(codeAttr.ExceptionTable(),
			me.class.constantPool)
	}
	me.exceptions = cfMethod.ExceptionsAttribute()
	me.annotationData = cfMethod.RuntimeVisibleAnnotationsAttributeData()
	me.parameterAnnotationData = cfMethod.RuntimeVisibleParameterAnnotationsAttributeData()
	me.annotationDefaultData = cfMethod.AnnotationDefaultAttributeData()
}

func (me *Method) calcArgSlotCount(paramTypes []string) {
	for _, paramType := range paramTypes {
		me.argSlotCount++
		if paramType == "J" || paramType == "D" {
			me.argSlotCount++
		}
	}
	if !me.IsStatic() {
		me.argSlotCount++ // `this` reference
	}
}

func (me *Method) injectCodeAttribute(returnType string) {
	me.maxStack = 4 // todo
	me.maxLocals = me.argSlotCount
	switch returnType[0] {
	case 'V':
		me.code = []byte{0xfe, 0xb1} // return
	case 'L', '[':
		me.code = []byte{0xfe, 0xb0} // areturn
	case 'D':
		me.code = []byte{0xfe, 0xaf} // dreturn
	case 'F':
		me.code = []byte{0xfe, 0xae} // freturn
	case 'J':
		me.code = []byte{0xfe, 0xad} // lreturn
	default:
		me.code = []byte{0xfe, 0xac} // ireturn
	}
}

func (me *Method) IsSynchronized() bool {
	return 0 != me.accessFlags&ACC_SYNCHRONIZED
}
func (me *Method) IsBridge() bool {
	return 0 != me.accessFlags&ACC_BRIDGE
}
func (me *Method) IsVarargs() bool {
	return 0 != me.accessFlags&ACC_VARARGS
}
func (me *Method) IsNative() bool {
	return 0 != me.accessFlags&ACC_NATIVE
}
func (me *Method) IsAbstract() bool {
	return 0 != me.accessFlags&ACC_ABSTRACT
}
func (me *Method) IsStrict() bool {
	return 0 != me.accessFlags&ACC_STRICT
}

// getters
func (me *Method) MaxStack() uint {
	return me.maxStack
}
func (me *Method) MaxLocals() uint {
	return me.maxLocals
}
func (me *Method) Code() []byte {
	return me.code
}
func (me *Method) ParameterAnnotationData() []byte {
	return me.parameterAnnotationData
}
func (me *Method) AnnotationDefaultData() []byte {
	return me.annotationDefaultData
}
func (me *Method) ParsedDescriptor() *MethodDescriptor {
	return me.parsedDescriptor
}
func (me *Method) ArgSlotCount() uint {
	return me.argSlotCount
}

func (me *Method) FindExceptionHandler(exClass *Class, pc int) int {
	handler := me.exceptionTable.findExceptionHandler(exClass, pc)
	if handler != nil {
		return handler.handlerPc
	}
	return -1
}

func (me *Method) GetLineNumber(pc int) int {
	if me.IsNative() {
		return -2
	}
	if me.lineNumberTable == nil {
		return -1
	}
	return me.lineNumberTable.GetLineNumber(pc)
}

func (me *Method) isConstructor() bool {
	return !me.IsStatic() && me.name == "<init>"
}
func (me *Method) isClinit() bool {
	return me.IsStatic() && me.name == "<clinit>"
}

// reflection
func (me *Method) ParameterTypes() []*Class {
	if me.argSlotCount == 0 {
		return nil
	}

	paramTypes := me.parsedDescriptor.parameterTypes
	paramClasses := make([]*Class, len(paramTypes))
	for i, paramType := range paramTypes {
		paramClassName := toClassName(paramType)
		paramClasses[i] = me.class.loader.LoadClass(paramClassName)
	}

	return paramClasses
}
func (me *Method) ReturnType() *Class {
	returnType := me.parsedDescriptor.returnType
	returnClassName := toClassName(returnType)
	return me.class.loader.LoadClass(returnClassName)
}
func (me *Method) ExceptionTypes() []*Class {
	if me.exceptions == nil {
		return nil
	}

	exIndexTable := me.exceptions.ExceptionIndexTable()
	exClasses := make([]*Class, len(exIndexTable))
	cp := me.class.constantPool

	for i, exIndex := range exIndexTable {
		classRef := cp.GetConstant(uint(exIndex)).(*ClassRef)
		exClasses[i] = classRef.ResolvedClass()
	}

	return exClasses
}
