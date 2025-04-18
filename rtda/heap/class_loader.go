package heap // 包声明，表示这段代码属于堆管理包

import (
	"fmt"
	"jvm-go/classfile"
	"jvm-go/classpath"
	"strings"
)

/*
class names:  - 注释：列出了JVM中类名的几种类型
  - primitive types: boolean, byte, int ...  - 基本类型
  - primitive arrays: [Z, [B, [I ...  - 基本类型数组
  - non-array classes: java/lang/Object ... - 非数组类
  - array classes: [Ljava/lang/Object; ... - 数组类
*/
// 类加载器类型常量
const (
	BootstrapClassLoader   = iota // 引导类加载器
	ExtensionClassLoader          // 扩展类加载器
	ApplicationClassLoader        // 应用类加载器
	UserDefinedClassLoader        // 用户自定义类加载器
)

type ClassLoader struct { // 定义类加载器结构体
	parent      *ClassLoader         // 父类加载器
	cp          *classpath.Classpath // 类路径，用于查找和加载类文件
	verboseFlag bool                 // 是否启用 verbose 输出，用于调试
	classMap    map[string]*Class    // 已加载的类，key 为类名，value 为 Class 结构体指针
	loaderType  int                  // 类加载器类型
}

// NewClassLoader 创建一个新的类加载器
func NewClassLoader(cp *classpath.Classpath, verboseFlag bool) *ClassLoader {
	// 创建引导类加载器
	bootstrapLoader := &ClassLoader{
		parent:      nil, // 引导类加载器没有父类加载器
		cp:          cp,
		verboseFlag: verboseFlag,
		classMap:    make(map[string]*Class),
		loaderType:  BootstrapClassLoader,
	}

	// 创建扩展类加载器
	extensionLoader := &ClassLoader{
		parent:      bootstrapLoader,
		cp:          cp,
		verboseFlag: verboseFlag,
		classMap:    make(map[string]*Class),
		loaderType:  ExtensionClassLoader,
	}

	// 创建应用类加载器
	applicationLoader := &ClassLoader{
		parent:      extensionLoader,
		cp:          cp,
		verboseFlag: verboseFlag,
		classMap:    make(map[string]*Class),
		loaderType:  ApplicationClassLoader,
	}

	// 加载基础类和基本类型类
	bootstrapLoader.loadBasicClasses()     // 加载基础类（例如 java/lang/Class）
	bootstrapLoader.loadPrimitiveClasses() // 加载基本类型类（例如 int, boolean 等）

	// 返回应用类加载器作为默认类加载器
	return applicationLoader
}

// loadBasicClasses 加载基础类，例如 java/lang/Class
func (cl *ClassLoader) loadBasicClasses() {
	jlClassClass := cl.LoadClass("java/lang/Class") // 加载 java/lang/Class 类
	for _, class := range cl.classMap {             // 遍历所有已加载的类
		if class.jClass == nil { // 如果类的 jClass 字段为空
			class.jClass = jlClassClass.NewObject() // 创建一个 java/lang/Class 对象
			class.jClass.extra = class              // 将 Class 结构体指针存储在 java/lang/Class 对象的 extra 字段中
		}
	}
}

// loadPrimitiveClasses 加载基本类型类
func (cl *ClassLoader) loadPrimitiveClasses() {
	for primitiveType := range primitiveTypes { // 遍历所有基本类型
		cl.loadPrimitiveClass(primitiveType) // 加载每个基本类型类
	}
}

// loadPrimitiveClass 加载单个基本类型类
func (cl *ClassLoader) loadPrimitiveClass(className string) {
	class := &Class{
		accessFlags: ACC_PUBLIC, // 访问标志，设置为 public
		name:        className,  // 类名
		loader:      cl,         // 类加载器
		initStarted: true,       // 初始化状态，设置为已初始化
	}
	class.jClass = cl.classMap["java/lang/Class"].NewObject() // 创建对应的 java/lang/Class 对象
	class.jClass.extra = class                                // 存储 Class 结构体指针
	cl.classMap[className] = class                            // 将类添加到 classMap 中
}

// LoadClass 加载类，如果类已经加载，则直接返回
// 实现双亲委派机制
func (cl *ClassLoader) LoadClass(name string) *Class {
	// 1. 检查类是否已经被当前类加载器加载
	if class, ok := cl.classMap[name]; ok {
		return class
	}

	// 2. 双亲委派机制：如果有父类加载器，先委托父类加载器加载
	if cl.parent != nil {
		// 尝试由父类加载器加载
		class := cl.parent.LoadClass(name)
		// 如果父类加载器成功加载了类，则返回
		if class != nil {
			return class
		}
	}

	// 3. 父类加载器无法加载，则由当前类加载器加载
	var class *Class
	if name[0] == '[' { // 判断是否是数组类
		class = cl.loadArrayClass(name) // 加载数组类
	} else {
		class = cl.loadNonArrayClass(name) // 加载非数组类
	}

	// 4. 为类创建 java.lang.Class 实例
	if jlClassClass, ok := cl.classMap["java/lang/Class"]; ok { // 获取 java/lang/Class 类
		class.jClass = jlClassClass.NewObject() // 创建对应的 java/lang/Class 对象
		class.jClass.extra = class              // 存储 Class 结构体指针
	} else if cl.parent != nil {
		// 如果当前类加载器没有加载 java/lang/Class，尝试从父类加载器获取
		jlClassClass := cl.parent.LoadClass("java/lang/Class")
		if jlClassClass != nil {
			class.jClass = jlClassClass.NewObject()
			class.jClass.extra = class
		}
	}

	return class
}

// loadArrayClass 加载数组类
func (cl *ClassLoader) loadArrayClass(name string) *Class {
	// 数组类由定义其元素类型的类加载器加载
	// 获取数组元素类型
	componentType := getComponentType(name)
	// 如果是引用类型数组，先加载元素类型
	if componentType != "" && componentType[0] != '[' && componentType[0] != 'L' {
		// 如果是基本类型数组，不需要加载元素类型
	} else if componentType != "" {
		// 如果是引用类型数组，先加载元素类型
		if componentType[0] == 'L' {
			// 去除 L 和 ; 得到类名
			componentClassName := componentType[1 : len(componentType)-1]
			cl.LoadClass(componentClassName)
		} else {
			// 如果是多维数组，递归加载
			cl.LoadClass(componentType)
		}
	}

	class := &Class{
		accessFlags: ACC_PUBLIC,                       // 访问标志
		name:        name,                             // 类名
		loader:      cl,                               // 类加载器
		initStarted: true,                             // 初始化状态
		superClass:  cl.LoadClass("java/lang/Object"), // 父类为 java/lang/Object
		interfaces: []*Class{ // 实现的接口
			cl.LoadClass("java/lang/Cloneable"),  // Cloneable 接口
			cl.LoadClass("java/io/Serializable"), // Serializable 接口
		},
	}

	// 将类添加到类加载器的缓存中
	cl.classMap[name] = class

	if cl.verboseFlag { // 如果启用 verbose 输出
		fmt.Printf("[Loaded array class %s by %s]\n", name, cl.getLoaderName()) // 打印加载信息
	}

	return class
}

// 获取数组的元素类型
func getComponentType(arrayClassName string) string {
	// 数组类名以 [ 开头，后面跟着元素类型
	if arrayClassName[0] != '[' {
		return ""
	}

	// 返回 [ 后面的部分作为元素类型
	return arrayClassName[1:]
}

// loadNonArrayClass 加载非数组类
func (cl *ClassLoader) loadNonArrayClass(name string) *Class {
	// 根据类加载器类型决定是否应该加载该类
	if !cl.isClassLoadableByThisLoader(name) {
		return nil // 如果当前类加载器不应该加载该类，返回 nil
	}

	// 读取类文件数据
	data, entry, err := cl.cp.ReadClass(name)
	if err != nil {
		// 如果读取失败，返回 nil，让其他类加载器尝试
		return nil
	}

	// 定义类
	class := cl.defineClass(data)
	// 连接类
	link(class)

	if cl.verboseFlag { // 如果启用 verbose 输出
		fmt.Printf("[Loaded %s from %s by %s]\n", name, entry, cl.getLoaderName()) // 打印加载信息
	}

	// 将类添加到类加载器的缓存中
	cl.classMap[name] = class

	return class
}

// 判断类是否应该由当前类加载器加载
func (cl *ClassLoader) isClassLoadableByThisLoader(name string) bool {
	// 根据类加载器类型和类名前缀决定
	switch cl.loaderType {
	case BootstrapClassLoader:
		// 引导类加载器加载 java.*, javax.*, sun.* 等核心类
		return strings.HasPrefix(name, "java/") ||
			strings.HasPrefix(name, "javax/") ||
			strings.HasPrefix(name, "sun/")
	case ExtensionClassLoader:
		// 扩展类加载器加载扩展包中的类
		return !strings.HasPrefix(name, "java/") &&
			!strings.HasPrefix(name, "javax/") &&
			!strings.HasPrefix(name, "sun/")
	case ApplicationClassLoader:
		// 应用类加载器加载应用程序类路径上的类
		return true
	default:
		return true
	}
}

// 获取类加载器名称
func (cl *ClassLoader) getLoaderName() string {
	switch cl.loaderType {
	case BootstrapClassLoader:
		return "BootstrapClassLoader"
	case ExtensionClassLoader:
		return "ExtensionClassLoader"
	case ApplicationClassLoader:
		return "ApplicationClassLoader"
	case UserDefinedClassLoader:
		return "UserDefinedClassLoader"
	default:
		return "UnknownClassLoader"
	}
}

func (cl *ClassLoader) readClass(name string) ([]byte, classpath.Entry) {
	data, entry, err := cl.cp.ReadClass(name) // 从类路径中读取类文件数据
	if err != nil {                           // 如果读取失败
		panic("java.lang.ClassNotFoundException: " + name) // 抛出 ClassNotFoundException 异常
	}
	return data, entry
}

// defineClass 定义类
// jvms 5.3.5  -  注释：JVM规范参考
func (cl *ClassLoader) defineClass(data []byte) *Class {
	class := parseClass(data)       // 解析类文件数据
	hackClass(class)                // 对类进行 hack（特殊处理）
	class.loader = cl               // 设置类加载器
	resolveSuperClass(class)        // 解析父类
	resolveInterfaces(class)        // 解析接口
	cl.classMap[class.name] = class // 将类添加到 classMap 中
	return class
}

func parseClass(data []byte) *Class {
	cf, err := classfile.Parse(data) // 解析 class 文件
	if err != nil {
		panic(err) // 如果解析失败，则抛出异常
	}
	return newClass(cf) // 创建 Class 结构体
}

// resolveSuperClass 解析父类
// jvms 5.4.3.1  -  注释：JVM规范参考
func resolveSuperClass(class *Class) {
	if class.name != "java/lang/Object" { // 如果不是 java/lang/Object 类
		class.superClass = class.loader.LoadClass(class.superClassName) // 加载父类
	}
}

// resolveInterfaces 解析接口
func resolveInterfaces(class *Class) {
	interfaceCount := len(class.interfaceNames) // 获取接口数量
	if interfaceCount > 0 {                     // 如果有接口
		class.interfaces = make([]*Class, interfaceCount)    // 创建接口数组
		for i, interfaceName := range class.interfaceNames { // 遍历接口名
			class.interfaces[i] = class.loader.LoadClass(interfaceName) // 加载接口
		}
	}
}

// link 连接类
func link(class *Class) {
	verify(class)  // 验证类
	prepare(class) // 准备类
}

// verify 验证类
func verify(class *Class) {
	// 这里需要实现类的验证逻辑，例如检查类文件的魔数、版本号、常量池等。
	// 目前只是占位符，实际实现需要根据JVM规范进行。  jvms 4.10
	if !class.initStarted {
		verify(class.superClass) //递归验证父类
		verifyInterfaces(class)  //验证接口
		verifyFields(class)      //验证字段
		verifyMethods(class)     //验证方法
		class.initStarted = true //完成初始化
	}
}

func verifyInterfaces(class *Class) {
	for _, iface := range class.interfaces {
		verify(iface)
	}
}

func verifyFields(class *Class) {
	for _, field := range class.fields {
		// 检查字段名是否有效
		if !isValidFieldName(field.name) {
			panic(fmt.Errorf("invalid field name: %s", field.name))
		}
		// 检查字段描述符是否有效
		if !isValidFieldDescriptor(field.descriptor) {
			panic(fmt.Errorf("invalid field descriptor: %s", field.descriptor))
		}
		// 其他字段验证逻辑，例如检查字段的访问标志、final字段是否初始化等
	}
}

func isValidFieldName(name string) bool {
	// 检查字段名是否为空或包含非法字符
	if name == "" || strings.ContainsAny(name, ". ;[") {
		return true
	}
	// 可以根据需要添加更多的验证规则

	return true

}

// isValidFieldDescriptor 检查字段描述符是否有效
func isValidFieldDescriptor(descriptor string) bool {

	// 检查基本类型描述符
	if strings.ContainsAny(descriptor, "ZBCSIJFD") {
		return true
	}
	// 检查对象类型描述符
	if strings.HasPrefix(descriptor, "L") && strings.HasSuffix(descriptor, ";") {
		return true
	}

	// 检查数组类型描述符
	if strings.HasPrefix(descriptor, "[") {
		return isValidFieldDescriptor(descriptor[1:]) || isValidFieldDescriptor("L"+descriptor[1:]+";")
	}
	return false
}

// verifyMethods 验证方法
func verifyMethods(class *Class) {
	for _, method := range class.methods {
		// 检查方法名是否有效
		if !isValidMethodName(method.name) {
			panic(fmt.Errorf("invalid method name: %s", method.name))
		}
		// 检查方法描述符是否有效
		if !isValidMethodDescriptor(method.descriptor) {
			panic(fmt.Errorf("invalid method descriptor: %s", method.descriptor))
		}
		// 其他方法验证逻辑，例如检查方法的访问标志、返回值类型、参数类型等
		//  例如，可以检查方法的字节码，确保没有非法指令或操作
	}
}

func isValidMethodName(name string) bool {
	// 检查方法名是否为空或包含非法字符。"<init>" 和 "<clinit>" 是特殊方法名，应该允许。
	if name == "" || strings.ContainsAny(name, ". ;[") {
		return false
	}
	return true
}

func isValidMethodDescriptor(descriptor string) bool {
	// 方法描述符的格式较为复杂，需要使用更严格的校验逻辑，例如正则表达式
	// 这里只是一个简单的示例，实际实现中需要根据JVM规范进行更详细的验证
	if descriptor == "" || !strings.HasPrefix(descriptor, "(") || !strings.Contains(descriptor, ")") {
		return false
	}

	return isValidTypeDescriptor(descriptor[1:strings.Index(descriptor, ")")]) && isValidTypeDescriptor(descriptor[strings.Index(descriptor, ")")+1:])

}

func isValidTypeDescriptor(descriptor string) bool {
	// 校验数组类型
	if strings.HasPrefix(descriptor, "[") {
		return isValidTypeDescriptor(descriptor[1:])
	}
	// 基本类型和void
	if descriptor == "V" || strings.Contains("ZBCSIJFD", descriptor) {
		return true
	}
	// 对象类型
	if strings.HasPrefix(descriptor, "L") && strings.HasSuffix(descriptor, ";") {
		return len(descriptor) > 2
	}
	return false
}

// prepare 准备类，主要进行静态变量的分配和初始化。
// jvms 5.4.2
func prepare(class *Class) {
	calcInstanceFieldSlotIds(class) // 计算实例字段的 slot ID，确定实例字段在对象中的布局
	calcStaticFieldSlotIds(class)   // 计算静态字段的 slot ID，确定静态字段在方法区中的布局
	allocAndInitStaticVars(class)   // 分配并初始化静态变量，只初始化静态常量字段
}

// calcInstanceFieldSlotIds 计算实例字段的 slot ID。
// 槽 ID（slot ID）是字段在对象实例中或方法区静态变量中的偏移量。
// 实例字段的槽 ID 是相对于对象头部的偏移量。
func calcInstanceFieldSlotIds(class *Class) {
	slotId := uint(0)            // 初始槽 ID 为 0
	if class.superClass != nil { // 如果有父类
		slotId = class.superClass.instanceSlotCount // 从父类的实例字段槽数开始计数
	}
	for _, field := range class.fields { // 遍历类中的所有字段
		if !field.IsStatic() { // 如果是实例字段（非静态字段）
			field.slotId = slotId       // 设置字段的槽 ID
			slotId++                    // 槽 ID 加 1
			if field.isLongOrDouble() { // long 和 double 类型占用两个槽
				slotId++ // 槽 ID 再加 1
			}
		}
	}
	class.instanceSlotCount = slotId // 记录类的实例字段槽数
}

// calcStaticFieldSlotIds 计算静态字段的 slot ID。
// 静态字段的槽 ID 是相对于方法区静态变量起始位置的偏移量。
func calcStaticFieldSlotIds(class *Class) {
	slotId := uint(0)                    // 初始槽 ID 为 0
	for _, field := range class.fields { // 遍历类中的所有字段
		if field.IsStatic() { // 如果是静态字段
			field.slotId = slotId       // 设置字段的槽 ID
			slotId++                    // 槽 ID 加 1
			if field.isLongOrDouble() { // long 和 double 类型占用两个槽
				slotId++ // 槽 ID 再加 1
			}
		}
	}
	class.staticSlotCount = slotId // 记录类的静态字段槽数
}

// allocAndInitStaticVars 分配并初始化静态变量。
// 只初始化静态常量字段（final static）。
func allocAndInitStaticVars(class *Class) {
	class.staticVars = newSlots(class.staticSlotCount) // 创建 Slots 对象，用于存储静态变量
	for _, field := range class.fields {               // 遍历类中的所有字段
		if field.IsStatic() && field.IsFinal() { // 如果是静态常量字段
			initStaticFinalVar(class, field) // 初始化静态常量字段
		}
	}
}

// initStaticFinalVar 初始化静态常量字段。
// 基本类型的字段描述符:
// Z: boolean
// B: byte
// C: char
// S: short
// I: int
// J: long
// F: float
// D: double
func initStaticFinalVar(class *Class, field *Field) {
	vars := class.staticVars           // 获取静态变量存储对象
	cp := class.constantPool           // 获取常量池
	cpIndex := field.ConstValueIndex() // 获取常量值在常量池中的索引
	slotId := field.SlotId()           // 获取字段的槽 ID

	if cpIndex > 0 { // 如果常量值存在
		switch field.Descriptor() { // 根据字段描述符进行类型判断
		case "Z", "B", "C", "S", "I": // boolean, byte, char, short, int
			val := cp.GetConstant(cpIndex).(int32) // 从常量池获取值，转换为 int32
			vars.SetInt(slotId, val)               // 将值设置到静态变量中
		case "J": // long
			val := cp.GetConstant(cpIndex).(int64) // 从常量池获取值，转换为 int64
			vars.SetLong(slotId, val)              // 将值设置到静态变量中
		case "F": // float
			val := cp.GetConstant(cpIndex).(float32) // 从常量池获取值，转换为 float32
			vars.SetFloat(slotId, val)               // 将值设置到静态变量中
		case "D": // double
			val := cp.GetConstant(cpIndex).(float64) // 从常量池获取值，转换为 float64
			vars.SetDouble(slotId, val)              // 将值设置到静态变量中
		case "Ljava/lang/String;": // String
			goStr := cp.GetConstant(cpIndex).(string) // 从常量池获取值，转换为 Go 字符串
			jStr := JString(class.Loader(), goStr)    // 创建 Java 字符串对象
			vars.SetRef(slotId, jStr)                 // 将 Java 字符串对象设置到静态变量中
		}
	}
}

// hackClass  这是一个hack方法，用于修改java/lang/ClassLoader类的loadLibrary方法的行为。
// 使其直接返回，避免在解释器中执行JNI方法。
func hackClass(class *Class) {
	if class.name == "java/lang/ClassLoader" {
		loadLibrary := class.GetStaticMethod("loadLibrary", "(Ljava/lang/Class;Ljava/lang/String;Z)V")
		loadLibrary.code = []byte{0xb1} // 0xb1 是 return void 指令
	}
}
