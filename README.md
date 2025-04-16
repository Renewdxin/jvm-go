# JVM-Go

一个用 Go 语言实现的 Java 虚拟机。本项目旨在深入理解 JVM 的工作原理，通过实现 JVM 规范来学习虚拟机设计。

## 项目概述

本项目实现了 JVM 的核心功能，包括：
- 类文件解析
- 类加载系统
- 运行时数据区
- 字节码解释器
- 部分本地方法

### 为什么用 Go 实现 JVM？

- Go 语言的并发特性和垃圾回收机制与 Java 类似
- Go 的语法简洁，标准库丰富，适合实现系统级程序
- 作为学习项目，用 Go 实现 JVM 可以深入理解两种语言的运行时机制

## 技术实现

### 1. 类文件格式 (Class File Format)

类文件解析主要在 `classfile` 包中实现：

```go
// 常量池实现
type ConstantPool struct {
    class  *Class
    consts []Constant
}

// 方法信息
type Method struct {
    ClassMember
    maxStack uint
    maxLocals uint
    code []byte
    // ...
}
```

### 2. 类加载系统 (Class Loading System)

类加载器实现在 `rtda/heap/class_loader.go`：

```go
type ClassLoader struct {
    cp          *classpath.Classpath
    classMap    map[string]*Class
}
```

支持：
- 类的加载、链接和初始化
- 数组类和普通类的区分处理
- 基本的类搜索机制

### 3. 运行时数据区 (Runtime Data Areas)

包含以下组件：

- JVM 栈和栈帧
```go
type Frame struct {
    localVars    LocalVars
    operandStack *OperandStack
    thread       *Thread
    method       *heap.Method
    nextPC       int
}
```

- 程序计数器
```go
type Thread struct {
    pc    int
    stack *Stack
}
```

- 堆和对象表示
```go
type Object struct {
    class *Class
    data  interface{}
}
```

### 4. 执行引擎 (Execution Engine)

采用字节码解释器实现，主要包括：

- 指令解释执行循环
- 方法调用处理
- 异常处理机制

### 5. 本地方法接口 (Native Interface)

实现了部分关键的本地方法：
- java.lang.Object
- java.lang.String
- java.lang.System
- java.io.FileInputStream/FileOutputStream

## 当前状态

### 已实现功能
- 基本的类加载和执行
- 大部分字节码指令的解释执行
- 简单的异常处理
- 基础的本地方法支持

### 限制和待实现功能
- 没有实现自己的垃圾回收器（依赖 Go 的 GC）
- 线程同步机制不完整
- 部分本地方法仍待实现
- JIT 编译器未实现

## 构建和运行

```bash
# 构建项目
go build

# 运行 Java class 文件
./jvm-go your-class-file.class
```

## 技术挑战

实现这个项目的主要技术挑战包括：
1. 正确实现 Java 的类加载机制
2. 准确解释执行 Java 字节码指令
3. 模拟 Java 运行时的内存管理和线程模型

## 学习资源

如果你想深入了解 JVM：
- [Java Virtual Machine Specification](https://docs.oracle.com/javase/specs/jvms/se8/html/index.html)
- [深入理解 Java 虚拟机](https://book.douban.com/subject/34907497/)

## 常见问题

### 问题：你在用 Go 实现 JVM 时，如何解析 class 文件？遇到了哪些挑战？

在实现 JVM 的 class 文件解析时，我主要采用了以下方法：

1. **字节流处理**：使用 Go 的 `encoding/binary` 包处理大端字节序（Big-Endian）的数据读取。JVM 规范要求 class 文件中的数据以大端格式存储，这与网络字节序一致。

```go
func (cr *ClassReader) readUint16() uint16 {
    val := binary.BigEndian.Uint16(cr.data) // 使用 BigEndian 字节序解码数据
    cr.data = cr.data[2:]                   // 后移两个字节
    return val
}
```

2. **魔数验证**：每个 class 文件开头都有魔数 `0xCAFEBABE`，用于验证文件格式。

```go
func (cf *ClassFile) readAndCheckMagic(reader *ClassReader) {
    magic := reader.readUint32()
    if magic != 0xCAFEBABE {
        panic("java.lang.ClassFormatError: magic!")
    }
}
```

3. **常量池解析**：常量池是 class 文件中最复杂的结构，包含多种类型的常量。我使用了接口和类型断言来处理不同类型的常量。

```go
type ConstantInfo interface {
    readInfo(reader *ClassReader)
}

func readConstantInfo(reader *ClassReader, cp ConstantPool) ConstantInfo {
    tag := reader.readUint8()
    c := newConstantInfo(tag, cp)
    c.readInfo(reader)
    return c
}
```

**主要挑战**：

1. **复杂的数据结构**：class 文件包含嵌套的数据结构，如常量池中的各种引用关系。
2. **字节序处理**：确保正确处理大端字节序的数据。
3. **特殊情况处理**：如 `CONSTANT_Long_info` 和 `CONSTANT_Double_info` 占用两个常量池项。
4. **属性表解析**：不同的属性有不同的格式，需要根据属性名动态解析。
5. **UTF-8 字符串处理**：Java 的 MUTF-8 编码与标准 UTF-8 有细微差别。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

[MIT License](LICENSE)
