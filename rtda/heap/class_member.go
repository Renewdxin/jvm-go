package heap

import "jvm-go/classfile"

type ClassMember struct {
	accessFlags    uint16
	name           string
	descriptor     string
	signature      string
	annotationData []byte // RuntimeVisibleAnnotations_attribute
	// class字段存放Class结构体指针，这样可以通过字段或方法访问到它所属的类
	class          *Class
}

func (clm *ClassMember) copyMemberInfo(memberInfo *classfile.MemberInfo) {
	clm.accessFlags = memberInfo.AccessFlags()
	clm.name = memberInfo.Name()
	clm.descriptor = memberInfo.Descriptor()
}

func (clm *ClassMember) IsPublic() bool {
	return 0 != clm.accessFlags&ACC_PUBLIC
}
func (clm *ClassMember) IsPrivate() bool {
	return 0 != clm.accessFlags&ACC_PRIVATE
}
func (clm *ClassMember) IsProtected() bool {
	return 0 != clm.accessFlags&ACC_PROTECTED
}
func (clm *ClassMember) IsStatic() bool {
	return 0 != clm.accessFlags&ACC_STATIC
}
func (clm *ClassMember) IsFinal() bool {
	return 0 != clm.accessFlags&ACC_FINAL
}
func (clm *ClassMember) IsSynthetic() bool {
	return 0 != clm.accessFlags&ACC_SYNTHETIC
}

// getters
func (clm *ClassMember) AccessFlags() uint16 {
	return clm.accessFlags
}
func (clm *ClassMember) Name() string {
	return clm.name
}
func (clm *ClassMember) Descriptor() string {
	return clm.descriptor
}
func (clm *ClassMember) Signature() string {
	return clm.signature
}
func (clm *ClassMember) AnnotationData() []byte {
	return clm.annotationData
}
func (clm *ClassMember) Class() *Class {
	return clm.class
}

// jvms 5.4.4
func (clm *ClassMember) isAccessibleTo(d *Class) bool {
	if clm.IsPublic() {
		return true
	}
	c := clm.class
	if clm.IsProtected() {
		return d == c || d.IsSubClassOf(c) ||
			c.GetPackageName() == d.GetPackageName()
	}
	if !clm.IsPrivate() {
		return c.GetPackageName() == d.GetPackageName()
	}
	return d == c
}
