package classfile

import "fmt"

/*
ClassFile {
    u4             magic;
    u2             minor_version;
    u2             major_version;
    u2             constant_pool_count;
    cp_info        constant_pool[constant_pool_count-1];
    u2             access_flags;
    u2             this_class;
    u2             super_class;
    u2             interfaces_count;
    u2             interfaces[interfaces_count];
    u2             fields_count;
    field_info     fields[fields_count];
    u2             methods_count;
    method_info    methods[methods_count];
    u2             attributes_count;
    attribute_info attributes[attributes_count];
}
*/

// ClassFile结构体包含类文件的所有信息
type ClassFile struct {
	//magic      uint32
	minorVersion uint16
	// 次版本号
	majorVersion uint16
	// 常量池
	constantPool ConstantPool
	// 访问标志
	accessFlags  uint16
	// 当前类索引
	thisClass    uint16
	// 父类索引
	superClass   uint16
	// 接口索引
	interfaces   []uint16
	// 字段
	fields       []*MemberInfo
	// 方法
	methods      []*MemberInfo
	// 属性
	attributes   []AttributeInfo
}

// Parse函数解析字节码，返回ClassFile结构体
func Parse(classData []byte) (cf *ClassFile, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	cr := &ClassReader{classData}
	cf = &ClassFile{}
	cf.read(cr)
	return
}

// read函数读取ClassFile结构体中的所有信息
func (cf *ClassFile) read(reader *ClassReader) {
	cf.readAndCheckMagic(reader)
	cf.readAndCheckVersion(reader)
	cf.constantPool = readConstantPool(reader)
	cf.accessFlags = reader.readUint16()
	cf.thisClass = reader.readUint16()
	cf.superClass = reader.readUint16()
	cf.interfaces = reader.readUint16s()
	cf.fields = readMembers(reader, cf.constantPool)
	cf.methods = readMembers(reader, cf.constantPool)
	cf.attributes = readAttributes(reader, cf.constantPool)
}

func (cf *ClassFile) readAndCheckMagic(reader *ClassReader) {
	magic := reader.readUint32()
	if magic != 0xCAFEBABE {
		panic("java.lang.ClassFormatError: magic!")
	}
}

func (cf *ClassFile) readAndCheckVersion(reader *ClassReader) {
	cf.minorVersion = reader.readUint16()
	cf.majorVersion = reader.readUint16()
	switch cf.majorVersion {
	case 45:
		return
	case 46, 47, 48, 49, 50, 51, 52:
		if cf.minorVersion == 0 {
			return
		}
	}

	panic("java.lang.UnsupportedClassVersionError!")
}

// Getter方法，暴露ClassFile结构体中的信息
func (cf *ClassFile) MinorVersion() uint16 {
	return cf.minorVersion
}
func (cf *ClassFile) MajorVersion() uint16 {
	return cf.majorVersion
}
func (cf *ClassFile) ConstantPool() ConstantPool {
	return cf.constantPool
}
func (cf *ClassFile) AccessFlags() uint16 {
	return cf.accessFlags
}
func (cf *ClassFile) Fields() []*MemberInfo {
	return cf.fields
}
func (cf *ClassFile) Methods() []*MemberInfo {
	return cf.methods
}

func (cf *ClassFile) ClassName() string {
	return cf.constantPool.getClassName(cf.thisClass)
}

// 常量池中获取超类名
func (cf *ClassFile) SuperClassName() string {
	if cf.superClass > 0 {
		return cf.constantPool.getClassName(cf.superClass)
	}
	return ""
}

func (cf *ClassFile) InterfaceNames() []string {
	interfaceNames := make([]string, len(cf.interfaces))
	for i, cpIndex := range cf.interfaces {
		interfaceNames[i] = cf.constantPool.getClassName(cpIndex)
	}
	return interfaceNames
}

func (cf *ClassFile) SourceFileAttribute() *SourceFileAttribute {
	for _, attrInfo := range cf.attributes {
		switch attrInfo.(type) {
		case *SourceFileAttribute:
			return attrInfo.(*SourceFileAttribute)
		}
	}
	return nil
}
