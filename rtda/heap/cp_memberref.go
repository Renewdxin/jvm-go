package heap

import "jvm-go/classfile"

// MemberRef 存放字段和方法符号引用共有的信息
type MemberRef struct {
	SymRef
	name       string
	descriptor string
}

// 从class文件内存储的字段或方法常量中提取数据
func (mr *MemberRef) copyMemberRefInfo(refInfo *classfile.ConstantMemberrefInfo) {
	mr.className = refInfo.ClassName()
	mr.name, mr.descriptor = refInfo.NameAndDescriptor()
}

func (mr *MemberRef) Name() string {
	return mr.name
}
func (mr *MemberRef) Descriptor() string {
	return mr.descriptor
}
