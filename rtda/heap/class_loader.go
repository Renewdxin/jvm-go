package heap

import "fmt"
import "jvm-go/classfile"
import "jvm-go/classpath"

/*
class names:
    - primitive types: boolean, byte, int ...
    - primitive arrays: [Z, [B, [I ...
    - non-array classes: java/lang/Object ...
    - array classes: [Ljava/lang/Object; ...
*/
type ClassLoader struct {
	cp          *classpath.Classpath
	verboseFlag bool
	classMap    map[string]*Class // loaded classes
}

func NewClassLoader(cp *classpath.Classpath, verboseFlag bool) *ClassLoader {
	loader := &ClassLoader{
		cp:          cp,
		verboseFlag: verboseFlag,
		classMap:    make(map[string]*Class),
	}

	loader.loadBasicClasses()
	loader.loadPrimitiveClasses()
	return loader
}

func (cl *ClassLoader) loadBasicClasses() {
	jlClassClass := cl.LoadClass("java/lang/Class")
	for _, class := range cl.classMap {
		if class.jClass == nil {
			class.jClass = jlClassClass.NewObject()
			class.jClass.extra = class
		}
	}
}

func (cl *ClassLoader) loadPrimitiveClasses() {
	for primitiveType, _ := range primitiveTypes {
		cl.loadPrimitiveClass(primitiveType)
	}
}

func (cl *ClassLoader) loadPrimitiveClass(className string) {
	class := &Class{
		accessFlags: ACC_PUBLIC, // todo
		name:        className,
		loader:      cl,
		initStarted: true,
	}
	class.jClass = cl.classMap["java/lang/Class"].NewObject()
	class.jClass.extra = class
	cl.classMap[className] = class
}

func (cl *ClassLoader) LoadClass(name string) *Class {
	if class, ok := cl.classMap[name]; ok {
		// already loaded
		return class
	}

	var class *Class
	if name[0] == '[' { // array class
		class = cl.loadArrayClass(name)
	} else {
		class = cl.loadNonArrayClass(name)
	}

	if jlClassClass, ok := cl.classMap["java/lang/Class"]; ok {
		class.jClass = jlClassClass.NewObject()
		class.jClass.extra = class
	}

	return class
}

func (cl *ClassLoader) loadArrayClass(name string) *Class {
	class := &Class{
		accessFlags: ACC_PUBLIC, // todo
		name:        name,
		loader:      cl,
		initStarted: true,
		superClass:  cl.LoadClass("java/lang/Object"),
		interfaces: []*Class{
			cl.LoadClass("java/lang/Cloneable"),
			cl.LoadClass("java/io/Serializable"),
		},
	}
	cl.classMap[name] = class
	return class
}

func (cl *ClassLoader) loadNonArrayClass(name string) *Class {
	data, entry := cl.readClass(name)
	class := cl.defineClass(data)
	link(class)

	if cl.verboseFlag {
		fmt.Printf("[Loaded %s from %s]\n", name, entry)
	}

	return class
}

func (cl *ClassLoader) readClass(name string) ([]byte, classpath.Entry) {
	data, entry, err := cl.cp.ReadClass(name)
	if err != nil {
		panic("java.lang.ClassNotFoundException: " + name)
	}
	return data, entry
}

// jvms 5.3.5
func (cl *ClassLoader) defineClass(data []byte) *Class {
	class := parseClass(data)
	hackClass(class)
	class.loader = cl
	resolveSuperClass(class)
	resolveInterfaces(class)
	cl.classMap[class.name] = class
	return class
}

func parseClass(data []byte) *Class {
	cf, err := classfile.Parse(data)
	if err != nil {
		//panic("java.lang.ClassFormatError")
		panic(err)
	}
	return newClass(cf)
}

// jvms 5.4.3.1
func resolveSuperClass(class *Class) {
	if class.name != "java/lang/Object" {
		class.superClass = class.loader.LoadClass(class.superClassName)
	}
}
func resolveInterfaces(class *Class) {
	interfaceCount := len(class.interfaceNames)
	if interfaceCount > 0 {
		class.interfaces = make([]*Class, interfaceCount)
		for i, interfaceName := range class.interfaceNames {
			class.interfaces[i] = class.loader.LoadClass(interfaceName)
		}
	}
}

func link(class *Class) {
	verify(class)
	prepare(class)
}

func verify(class *Class) {
	// todo
}

// jvms 5.4.2
func prepare(class *Class) {
	calcInstanceFieldSlotIds(class)
	calcStaticFieldSlotIds(class)
	allocAndInitStaticVars(class)
}

func calcInstanceFieldSlotIds(class *Class) {
	slotId := uint(0)
	if class.superClass != nil {
		slotId = class.superClass.instanceSlotCount
	}
	for _, field := range class.fields {
		if !field.IsStatic() {
			field.slotId = slotId
			slotId++
			if field.isLongOrDouble() {
				slotId++
			}
		}
	}
	class.instanceSlotCount = slotId
}

func calcStaticFieldSlotIds(class *Class) {
	slotId := uint(0)
	for _, field := range class.fields {
		if field.IsStatic() {
			field.slotId = slotId
			slotId++
			if field.isLongOrDouble() {
				slotId++
			}
		}
	}
	class.staticSlotCount = slotId
}

func allocAndInitStaticVars(class *Class) {
	class.staticVars = newSlots(class.staticSlotCount)
	for _, field := range class.fields {
		if field.IsStatic() && field.IsFinal() {
			initStaticFinalVar(class, field)
		}
	}
}

func initStaticFinalVar(class *Class, field *Field) {
	vars := class.staticVars
	cp := class.constantPool
	cpIndex := field.ConstValueIndex()
	slotId := field.SlotId()

	if cpIndex > 0 {
		switch field.Descriptor() {
		case "Z", "B", "C", "S", "I":
			val := cp.GetConstant(cpIndex).(int32)
			vars.SetInt(slotId, val)
		case "J":
			val := cp.GetConstant(cpIndex).(int64)
			vars.SetLong(slotId, val)
		case "F":
			val := cp.GetConstant(cpIndex).(float32)
			vars.SetFloat(slotId, val)
		case "D":
			val := cp.GetConstant(cpIndex).(float64)
			vars.SetDouble(slotId, val)
		case "Ljava/lang/String;":
			goStr := cp.GetConstant(cpIndex).(string)
			jStr := JString(class.Loader(), goStr)
			vars.SetRef(slotId, jStr)
		}
	}
}

// todo
func hackClass(class *Class) {
	if class.name == "java/lang/ClassLoader" {
		loadLibrary := class.GetStaticMethod("loadLibrary", "(Ljava/lang/Class;Ljava/lang/String;Z)V")
		loadLibrary.code = []byte{0xb1} // return void
	}
}
