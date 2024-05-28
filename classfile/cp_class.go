package classfile

type ConstantClassInfo struct {
	cp                ConstantPool
	// 常量池索引
	nameIndex       uint16
}

func (self *ConstantClassInfo) readInfo(reader *ClassReader) {
	self.nameIndex = reader.readUint16()
}
func (self *ConstantClassInfo) Name() string {
	return self.cp.getUtf8(self.nameIndex)
}