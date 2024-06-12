package classfile

/*
CONSTANT_Class_info {
    u1 tag;
    u2 name_index;
}
*/
type ConstantClassInfo struct {
	cp        ConstantPool
	nameIndex uint16
}

func (cci *ConstantClassInfo) readInfo(reader *ClassReader) {
	cci.nameIndex = reader.readUint16()
}
func (cci *ConstantClassInfo) Name() string {
	return cci.cp.getUtf8(cci.nameIndex)
}
