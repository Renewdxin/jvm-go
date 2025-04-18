# Java类加载器的双亲委派机制实现

## 1. 双亲委派机制概述

双亲委派机制是Java类加载器的一个重要特性，它的基本原理是：

1. 当一个类加载器收到类加载请求时，它首先将这个请求委派给父类加载器去完成
2. 只有当父类加载器无法找到这个类时（在其搜索范围内没有找到所需的类），子类加载器才会尝试自己去加载
3. 这种机制可以确保Java核心类库的类型安全，防止用户自定义的类替换Java核心类库的类

## 2. 类加载器层次结构

在标准的Java虚拟机中，类加载器形成了一个层次结构：

1. **引导类加载器（Bootstrap ClassLoader）**：
   - 负责加载Java核心类库，如`java.*`、`javax.*`、`sun.*`等包中的类
   - 通常由C++实现，是虚拟机的一部分

2. **扩展类加载器（Extension ClassLoader）**：
   - 负责加载Java扩展类库，如`jre/lib/ext`目录下的类
   - 是引导类加载器的子类加载器

3. **应用类加载器（Application ClassLoader）**：
   - 负责加载应用程序类路径（classpath）上的类
   - 是扩展类加载器的子类加载器
   - 也称为系统类加载器（System ClassLoader）

4. **用户自定义类加载器（User-Defined ClassLoader）**：
   - 由用户自定义的类加载器，通常继承自`java.lang.ClassLoader`
   - 是应用类加载器的子类加载器

## 3. 实现思路

### 3.1 类加载器结构设计

首先，我们需要修改`ClassLoader`结构体，添加父类加载器引用和类加载器类型：

```go
// 类加载器类型常量
const (
    BootstrapClassLoader = iota // 引导类加载器
    ExtensionClassLoader         // 扩展类加载器
    ApplicationClassLoader       // 应用类加载器
    UserDefinedClassLoader       // 用户自定义类加载器
)

type ClassLoader struct {
    parent      *ClassLoader           // 父类加载器
    cp          *classpath.Classpath   // 类路径
    verboseFlag bool                   // 是否启用 verbose 输出
    classMap    map[string]*Class      // 已加载的类缓存
    loaderType  int                    // 类加载器类型
}
```

### 3.2 创建类加载器层次结构

在`NewClassLoader`函数中，我们创建了类加载器的层次结构：

```go
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
    bootstrapLoader.loadBasicClasses()
    bootstrapLoader.loadPrimitiveClasses()

    // 返回应用类加载器作为默认类加载器
    return applicationLoader
}
```

### 3.3 实现双亲委派机制

在`LoadClass`方法中，我们实现了双亲委派机制的核心逻辑：

```go
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
    if jlClassClass, ok := cl.classMap["java/lang/Class"]; ok {
        class.jClass = jlClassClass.NewObject()
        class.jClass.extra = class
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
```

### 3.4 类加载器职责划分

我们通过`isClassLoadableByThisLoader`方法判断一个类是否应该由当前类加载器加载：

```go
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
```

### 3.5 数组类加载

数组类的加载需要特殊处理，因为数组类是由JVM在运行时动态创建的，而不是从类文件加载的：

```go
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
            componentClassName := componentType[1:len(componentType)-1]
            cl.LoadClass(componentClassName)
        } else {
            // 如果是多维数组，递归加载
            cl.LoadClass(componentType)
        }
    }

    class := &Class{
        accessFlags: ACC_PUBLIC,
        name:        name,
        loader:      cl,
        initStarted: true,
        superClass:  cl.LoadClass("java/lang/Object"),
        interfaces: []*Class{
            cl.LoadClass("java/lang/Cloneable"),
            cl.LoadClass("java/io/Serializable"),
        },
    }

    // 将类添加到类加载器的缓存中
    cl.classMap[name] = class

    if cl.verboseFlag {
        fmt.Printf("[Loaded array class %s by %s]\n", name, cl.getLoaderName())
    }

    return class
}
```

## 4. 双亲委派机制的优点

1. **确保类型安全**：
   - 防止用户自定义的类替换Java核心类库的类
   - 例如，用户无法定义自己的`java.lang.Object`类来替换JDK中的`Object`类

2. **避免类的重复加载**：
   - 父类加载器已经加载的类，子类加载器不会再次加载
   - 节省内存空间，提高性能

3. **保证核心类库的一致性**：
   - 确保Java核心类库的类只会被引导类加载器加载
   - 维护Java平台的完整性和一致性

## 5. 实现中的关键点

1. **类加载器缓存**：
   - 每个类加载器都有自己的类缓存（`classMap`）
   - 加载成功后，将类添加到缓存中，避免重复加载

2. **类加载器层次结构**：
   - 明确定义了类加载器的层次关系
   - 每个类加载器都有一个父类加载器（除了引导类加载器）

3. **类加载的职责划分**：
   - 根据类的包名前缀，确定由哪个类加载器负责加载
   - 核心类库由引导类加载器加载，扩展类库由扩展类加载器加载，应用程序类由应用类加载器加载

4. **数组类的特殊处理**：
   - 数组类由定义其元素类型的类加载器创建
   - 对于多维数组，需要递归处理

## 6. 总结

通过实现双亲委派机制，我们的JVM现在能够更好地保证类型安全，并且符合Java类加载的规范。这种机制确保了Java核心类库的类型安全，防止用户自定义的类替换Java核心类库的类，同时也提高了类加载的效率，避免了重复加载。

双亲委派模型虽然不是强制性的，但它是Java类加载器的一个重要设计原则，有助于维护Java平台的安全性和一致性。在某些特殊情况下，可能需要打破双亲委派模型，例如JNDI、JDBC等服务提供者接口（SPI）的实现，但这需要谨慎处理，以避免安全风险。
