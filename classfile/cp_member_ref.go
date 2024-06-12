package classfile

/*
CONSTANT_Fieldref_info {
    u1 tag;
    u2 class_index;
    u2 name_and_type_index;
}
CONSTANT_Methodref_info {
    u1 tag;
    u2 class_index;
    u2 name_and_type_index;
}
CONSTANT_InterfaceMethodref_info {
    u1 tag;
    u2 class_index;
    u2 name_and_type_index;
}
*/
type ConstantFieldrefInfo struct{ ConstantMemberrefInfo }
type ConstantMethodrefInfo struct{ ConstantMemberrefInfo }
type ConstantInterfaceMethodrefInfo struct{ ConstantMemberrefInfo }

type ConstantMemberrefInfo struct {
	cp               ConstantPool
	classIndex       uint16
	nameAndTypeIndex uint16
}

func (cni *ConstantMemberrefInfo) readInfo(reader *ClassReader) {
	cni.classIndex = reader.readUint16()
	cni.nameAndTypeIndex = reader.readUint16()
}

func (cni *ConstantMemberrefInfo) ClassName() string {
	return cni.cp.getClassName(cni.classIndex)
}
func (cni *ConstantMemberrefInfo) NameAndDescriptor() (string, string) {
	return cni.cp.getNameAndType(cni.nameAndTypeIndex)
}
