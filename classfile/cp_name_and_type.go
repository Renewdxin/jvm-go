package classfile

/*
	CONSTANT_NameAndType_info {
	    u1 tag;
	    u2 name_index;
	    u2 descriptor_index;
	}
*/
type ConstantNameAndTypeInfo struct {
	nameIndex       uint16
	descriptorIndex uint16
}

func (cnti *ConstantNameAndTypeInfo) readInfo(reader *ClassReader) {
	cnti.nameIndex = reader.readUint16()
	cnti.descriptorIndex = reader.readUint16()
}
