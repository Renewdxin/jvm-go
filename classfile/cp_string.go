package classfile

/*
CONSTANT_String_info {
    u1 tag;
    u2 string_index;
}
*/
type ConstantStringInfo struct {
	cp          ConstantPool
	stringIndex uint16
}

func (csi *ConstantStringInfo) readInfo(reader *ClassReader) {
	csi.stringIndex = reader.readUint16()
}
func (csi *ConstantStringInfo) String() string {
	return csi.cp.getUtf8(csi.stringIndex)
}
