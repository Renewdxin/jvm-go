package heap

import "strings"
import "jvm-go/classfile"

// name, superClassName and interfaceNames are all binary names(jvms8-4.2.1)
type Class struct {
	accessFlags       uint16
	name              string // thisClassName
	superClassName    string
	interfaceNames    []string
	constantPool      *ConstantPool
	fields            []*Field
	methods           []*Method
	sourceFile        string
	loader            *ClassLoader
	superClass        *Class
	interfaces        []*Class
	instanceSlotCount uint
	staticSlotCount   uint
	staticVars        Slots
	initStarted       bool
	jClass            *Object
}

func newClass(cf *classfile.ClassFile) *Class {
	class := &Class{}
	class.accessFlags = cf.AccessFlags()
	class.name = cf.ClassName()
	class.superClassName = cf.SuperClassName()
	class.interfaceNames = cf.InterfaceNames()
	class.constantPool = newConstantPool(class, cf.ConstantPool())
	class.fields = newFields(class, cf.Fields())
	class.methods = newMethods(class, cf.Methods())
	class.sourceFile = getSourceFile(cf)
	return class
}

func getSourceFile(cf *classfile.ClassFile) string {
	if sfAttr := cf.SourceFileAttribute(); sfAttr != nil {
		return sfAttr.FileName()
	}
	return "Unknown" // todo
}

func (cl *Class) IsPublic() bool {
	return 0 != cl.accessFlags&ACC_PUBLIC
}
func (cl *Class) IsFinal() bool {
	return 0 != cl.accessFlags&ACC_FINAL
}
func (cl *Class) IsSuper() bool {
	return 0 != cl.accessFlags&ACC_SUPER
}
func (cl *Class) IsInterface() bool {
	return 0 != cl.accessFlags&ACC_INTERFACE
}
func (cl *Class) IsAbstract() bool {
	return 0 != cl.accessFlags&ACC_ABSTRACT
}
func (cl *Class) IsSynthetic() bool {
	return 0 != cl.accessFlags&ACC_SYNTHETIC
}
func (cl *Class) IsAnnotation() bool {
	return 0 != cl.accessFlags&ACC_ANNOTATION
}
func (cl *Class) IsEnum() bool {
	return 0 != cl.accessFlags&ACC_ENUM
}

// getters
func (cl *Class) AccessFlags() uint16 {
	return cl.accessFlags
}
func (cl *Class) Name() string {
	return cl.name
}
func (cl *Class) ConstantPool() *ConstantPool {
	return cl.constantPool
}
func (cl *Class) Fields() []*Field {
	return cl.fields
}
func (cl *Class) Methods() []*Method {
	return cl.methods
}
func (cl *Class) SourceFile() string {
	return cl.sourceFile
}
func (cl *Class) Loader() *ClassLoader {
	return cl.loader
}
func (cl *Class) SuperClass() *Class {
	return cl.superClass
}
func (cl *Class) Interfaces() []*Class {
	return cl.interfaces
}
func (cl *Class) StaticVars() Slots {
	return cl.staticVars
}
func (cl *Class) InitStarted() bool {
	return cl.initStarted
}
func (cl *Class) JClass() *Object {
	return cl.jClass
}

func (cl *Class) StartInit() {
	cl.initStarted = true
}

// jvms 5.4.4
func (cl *Class) isAccessibleTo(other *Class) bool {
	return cl.IsPublic() ||
		cl.GetPackageName() == other.GetPackageName()
}

func (cl *Class) GetPackageName() string {
	if i := strings.LastIndex(cl.name, "/"); i >= 0 {
		return cl.name[:i]
	}
	return ""
}

func (cl *Class) GetMainMethod() *Method {
	return cl.getMethod("main", "([Ljava/lang/String;)V", true)
}
func (cl *Class) GetClinitMethod() *Method {
	return cl.getMethod("<clinit>", "()V", true)
}

func (cl *Class) getMethod(name, descriptor string, isStatic bool) *Method {
	for c := cl; c != nil; c = c.superClass {
		for _, method := range c.methods {
			if method.IsStatic() == isStatic &&
				method.name == name &&
				method.descriptor == descriptor {

				return method
			}
		}
	}
	return nil
}

func (cl *Class) getField(name, descriptor string, isStatic bool) *Field {
	for c := cl; c != nil; c = c.superClass {
		for _, field := range c.fields {
			if field.IsStatic() == isStatic &&
				field.name == name &&
				field.descriptor == descriptor {

				return field
			}
		}
	}
	return nil
}

func (cl *Class) isJlObject() bool {
	return cl.name == "java/lang/Object"
}
func (cl *Class) isJlCloneable() bool {
	return cl.name == "java/lang/Cloneable"
}
func (cl *Class) isJioSerializable() bool {
	return cl.name == "java/io/Serializable"
}

func (cl *Class) NewObject() *Object {
	return newObject(cl)
}

func (cl *Class) ArrayClass() *Class {
	arrayClassName := getArrayClassName(cl.name)
	return cl.loader.LoadClass(arrayClassName)
}

func (cl *Class) JavaName() string {
	return strings.Replace(cl.name, "/", ".", -1)
}

func (cl *Class) IsPrimitive() bool {
	_, ok := primitiveTypes[cl.name]
	return ok
}

func (cl *Class) GetInstanceMethod(name, descriptor string) *Method {
	return cl.getMethod(name, descriptor, false)
}
func (cl *Class) GetStaticMethod(name, descriptor string) *Method {
	return cl.getMethod(name, descriptor, true)
}

// reflection
func (cl *Class) GetRefVar(fieldName, fieldDescriptor string) *Object {
	field := cl.getField(fieldName, fieldDescriptor, true)
	return cl.staticVars.GetRef(field.slotId)
}
func (cl *Class) SetRefVar(fieldName, fieldDescriptor string, ref *Object) {
	field := cl.getField(fieldName, fieldDescriptor, true)
	cl.staticVars.SetRef(field.slotId, ref)
}

func (cl *Class) GetFields(publicOnly bool) []*Field {
	if publicOnly {
		publicFields := make([]*Field, 0, len(cl.fields))
		for _, field := range cl.fields {
			if field.IsPublic() {
				publicFields = append(publicFields, field)
			}
		}
		return publicFields
	} else {
		return cl.fields
	}
}

func (cl *Class) GetConstructor(descriptor string) *Method {
	return cl.GetInstanceMethod("<init>", descriptor)
}

func (cl *Class) GetConstructors(publicOnly bool) []*Method {
	constructors := make([]*Method, 0, len(cl.methods))
	for _, method := range cl.methods {
		if method.isConstructor() {
			if !publicOnly || method.IsPublic() {
				constructors = append(constructors, method)
			}
		}
	}
	return constructors
}

func (cl *Class) GetMethods(publicOnly bool) []*Method {
	methods := make([]*Method, 0, len(cl.methods))
	for _, method := range cl.methods {
		if !method.isClinit() && !method.isConstructor() {
			if !publicOnly || method.IsPublic() {
				methods = append(methods, method)
			}
		}
	}
	return methods
}
