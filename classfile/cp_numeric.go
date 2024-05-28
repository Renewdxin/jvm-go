package classfile

import "math"

type ConstantIntegerInfo struct {
	val int32
}

func (self *ConstantIntegerInfo) readInfo(reader *ClassReader) {
	bytes := reader.readUint32()
	self.val = int32(bytes)
}

type ConstantFloatInfo struct {
	val float32
}
func (self *ConstantFloatInfo) readInfo(reader *ClassReader) {
	bytes := reader.readUint32()
	self.val = math.Float32frombits(bytes)

	
type ConstantLongInfo struct {
	val int64
}

// 读取一个uint64数据，然后把它转型成int64类型
func (self *ConstantLongInfo) readInfo(reader *ClassReader) {
	bytes := reader.readUint64()
	self.val = int64(bytes)
}

type ConstantDoubleInfo struct {
	val float64
}

// 读取一个uint64数据，然后调用math包的Float64frombits()函数把它转换成float64类型
func (self *ConstantDoubleInfo) readInfo(reader *ClassReader) {
	bytes := reader.readUint64()
	self.val = math.Float64frombits(bytes)
}

