package heap

import "fmt"
import "jvm-go/classfile"

type Constant interface{}

/*
	运行时常量池主要存放两类信息：字面量（literal）和符号引用（symbolic reference）。
	字面量包括整数、浮点数和字符串字面量；
	符号引用包括类符号引用、字段符号引用、方法符号引用和接口方法符号引用。
*/

type ConstantPool struct {
	class  *Class
	consts []Constant
}

func newConstantPool(class *Class, cfCp classfile.ConstantPool) *ConstantPool {
	cpCount := len(cfCp)
	consts := make([]Constant, cpCount)
	rtCp := &ConstantPool{class, consts}

	for i := 1; i < cpCount; i++ {
		cpInfo := cfCp[i]

		switch cpInfo.(type) {
		case *classfile.ConstantIntegerInfo:
			intInfo := cpInfo.(*classfile.ConstantIntegerInfo)
			consts[i] = intInfo.Value()
		case *classfile.ConstantFloatInfo:
			floatInfo := cpInfo.(*classfile.ConstantFloatInfo)
			consts[i] = floatInfo.Value()
		case *classfile.ConstantLongInfo:
			longInfo := cpInfo.(*classfile.ConstantLongInfo)
			consts[i] = longInfo.Value()
			i++
		case *classfile.ConstantDoubleInfo:
			doubleInfo := cpInfo.(*classfile.ConstantDoubleInfo)
			consts[i] = doubleInfo.Value()
			i++
		case *classfile.ConstantStringInfo:
			stringInfo := cpInfo.(*classfile.ConstantStringInfo)
			consts[i] = stringInfo.String()
		case *classfile.ConstantClassInfo:
			classInfo := cpInfo.(*classfile.ConstantClassInfo)
			consts[i] = newClassRef(rtCp, classInfo)
		case *classfile.ConstantFieldrefInfo:
			fieldrefInfo := cpInfo.(*classfile.ConstantFieldrefInfo)
			consts[i] = newFieldRef(rtCp, fieldrefInfo)
		case *classfile.ConstantMethodrefInfo:
			methodrefInfo := cpInfo.(*classfile.ConstantMethodrefInfo)
			consts[i] = newMethodRef(rtCp, methodrefInfo)
		case *classfile.ConstantInterfaceMethodrefInfo:
			methodrefInfo := cpInfo.(*classfile.ConstantInterfaceMethodrefInfo)
			consts[i] = newInterfaceMethodRef(rtCp, methodrefInfo)
		default:
			// todo
			consts[i] = nil
		}
	}
	return rtCp
}

func (cp *ConstantPool) GetConstant(index uint) Constant {
	if c := cp.consts[index]; c != nil {
		return c
	}
	panic(fmt.Sprintf("No constants at index %d", index))
}
