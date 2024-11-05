package heap

import "jvm-go/classfile"

// Method 结构体表示一个方法
type Method struct {
	ClassMember // 继承自 ClassMember，包含 accessFlags, name, descriptor, attributes 等字段
	// 存放操作数栈的最大深度
	maxStack uint
	// 存放局部变量表的大小
	maxLocals uint
	// 存放方法的字节码
	code []byte
	// 存放异常处理表
	exceptionTable ExceptionTable // todo: 重命名为更具描述性的名称
	// 存放行号表，用于调试
	// 程序抛出异常或使用调试器单步执行时，虚拟机可以利用行号表显示程序当前执行到源代码的哪一行
	lineNumberTable *classfile.LineNumberTableAttribute
	// 存放方法抛出的异常类型
	exceptions *classfile.ExceptionsAttribute // todo: 重命名为更具描述性的名称
	// 存放参数的注解数据
	parameterAnnotationData []byte // RuntimeVisibleParameterAnnotations_attribute
	// 存放注解的默认值
	annotationDefaultData []byte // AnnotationDefault_attribute
	// 存放解析后的方法描述符
	parsedDescriptor *MethodDescriptor
	// 存放参数占用的局部变量槽数量，用于确定参数在局部变量表中的起始位置
	argSlotCount uint
}

// newMethods 函数根据 class 文件中的方法信息创建 Method 对象数组
func newMethods(class *Class, cfMethods []*classfile.MemberInfo) []*Method {
	methods := make([]*Method, len(cfMethods))
	for i, cfMethod := range cfMethods {
		methods[i] = newMethod(class, cfMethod)
	}
	return methods
}

// newMethod 函数根据 class 文件中的方法信息创建一个 Method 对象
func newMethod(class *Class, cfMethod *classfile.MemberInfo) *Method {
	method := &Method{}
	method.class = class
	method.copyMemberInfo(cfMethod)                // 拷贝方法的基本信息，例如访问标志、名称、描述符等
	method.copyAttributes(cfMethod)                // 拷贝方法的属性，例如代码属性、异常属性等
	md := parseMethodDescriptor(method.descriptor) // 解析方法描述符
	method.parsedDescriptor = md
	method.calcArgSlotCount(md.parameterTypes) // 计算参数占用的局部变量槽数量
	if method.IsNative() {                     // 如果是 native 方法，则注入代码属性
		method.injectCodeAttribute(md.returnType)
	}
	return method
}

// copyAttributes 函数从 class 文件中拷贝方法的属性
func (me *Method) copyAttributes(cfMethod *classfile.MemberInfo) {
	if codeAttr := cfMethod.CodeAttribute(); codeAttr != nil { // 获取 Code 属性
		me.maxStack = codeAttr.MaxStack()                        // 获取操作数栈的最大深度
		me.maxLocals = codeAttr.MaxLocals()                      // 获取局部变量表的大小
		me.code = codeAttr.Code()                                // 获取字节码
		me.lineNumberTable = codeAttr.LineNumberTableAttribute() // 获取行号表
		me.exceptionTable = newExceptionTable(codeAttr.ExceptionTable(),
			me.class.constantPool) // 创建异常处理表
	}
	me.exceptions = cfMethod.ExceptionsAttribute()                                          // 获取方法抛出的异常类型
	me.annotationData = cfMethod.RuntimeVisibleAnnotationsAttributeData()                   // 获取方法的注解数据  //此处代码中没有这个字段，应该是笔误
	me.parameterAnnotationData = cfMethod.RuntimeVisibleParameterAnnotationsAttributeData() // 获取参数的注解数据
	me.annotationDefaultData = cfMethod.AnnotationDefaultAttributeData()                    // 获取注解的默认值
}

// calcArgSlotCount 函数计算参数占用的局部变量槽数量
func (me *Method) calcArgSlotCount(paramTypes []string) {
	for _, paramType := range paramTypes {
		me.argSlotCount++
		if paramType == "J" || paramType == "D" { // long 和 double 类型占用两个局部变量槽
			me.argSlotCount++
		}
	}
	if !me.IsStatic() { // 非静态方法需要额外一个局部变量槽来存储 this 引用
		me.argSlotCount++ // `this` 引用
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
