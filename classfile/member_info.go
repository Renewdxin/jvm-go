package classfile

/*
field_info {
    u2             access_flags;
    u2             name_index;
    u2             descriptor_index;
    u2             attributes_count;
    attribute_info attributes[attributes_count];
}
method_info {
    u2             access_flags;
    u2             name_index;
    u2             descriptor_index;
    u2             attributes_count;
    attribute_info attributes[attributes_count];
}
*/

type MemberInfo struct {
	cp              ConstantPool
	accessFlags     uint16
	nameIndex       uint16
	descriptorIndex uint16
	attributes      []AttributeInfo
}

// read field or method table
func readMembers(reader *ClassReader, cp ConstantPool) []*MemberInfo {
	memberCount := reader.readUint16()
	members := make([]*MemberInfo, memberCount)
	for i := range members {
		members[i] = readMember(reader, cp)
	}
	return members
}

func readMember(reader *ClassReader, cp ConstantPool) *MemberInfo {
	return &MemberInfo{
		cp:              cp,
		accessFlags:     reader.readUint16(),
		nameIndex:       reader.readUint16(),
		descriptorIndex: reader.readUint16(),
		attributes:      readAttributes(reader, cp),
	}
}

func (mIn *MemberInfo) AccessFlags() uint16 {
	return mIn.accessFlags
}
func (mIn *MemberInfo) Name() string {
	return mIn.cp.getUtf8(mIn.nameIndex)
}
func (mIn *MemberInfo) Descriptor() string {
	return mIn.cp.getUtf8(mIn.descriptorIndex)
}

func (mIn *MemberInfo) CodeAttribute() *CodeAttribute {
	for _, attrInfo := range mIn.attributes {
		switch attrInfo.(type) {
		case *CodeAttribute:
			return attrInfo.(*CodeAttribute)
		}
	}
	return nil
}

func (mIn *MemberInfo) ConstantValueAttribute() *ConstantValueAttribute {
	for _, attrInfo := range mIn.attributes {
		switch attrInfo.(type) {
		case *ConstantValueAttribute:
			return attrInfo.(*ConstantValueAttribute)
		}
	}
	return nil
}

func (mIn *MemberInfo) ExceptionsAttribute() *ExceptionsAttribute {
	for _, attrInfo := range mIn.attributes {
		switch attrInfo.(type) {
		case *ExceptionsAttribute:
			return attrInfo.(*ExceptionsAttribute)
		}
	}
	return nil
}

func (mIn *MemberInfo) RuntimeVisibleAnnotationsAttributeData() []byte {
	return mIn.getUnparsedAttributeData("RuntimeVisibleAnnotations")
}
func (mIn *MemberInfo) RuntimeVisibleParameterAnnotationsAttributeData() []byte {
	return mIn.getUnparsedAttributeData("RuntimeVisibleParameterAnnotationsAttribute")
}
func (mIn *MemberInfo) AnnotationDefaultAttributeData() []byte {
	return mIn.getUnparsedAttributeData("AnnotationDefault")
}

func (mIn *MemberInfo) getUnparsedAttributeData(name string) []byte {
	for _, attrInfo := range mIn.attributes {
		switch attrInfo.(type) {
		case *UnparsedAttribute:
			unparsedAttr := attrInfo.(*UnparsedAttribute)
			if unparsedAttr.name == name {
				return unparsedAttr.info
			}
		}
	}
	return nil
}
