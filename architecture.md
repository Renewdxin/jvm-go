# JVM-Go 项目架构

## 整体架构图

下图展示了 JVM-Go 项目的整体架构，包括从启动到执行 Java 程序的完整流程。

```mermaid
graph TD
    A[main.go] --> B[解析命令行参数]
    B --> C[创建JVM实例]
    C --> D[启动JVM]
    D --> E[初始化VM]
    D --> F[执行main方法]

    E --> G[加载sun/misc/VM类]
    E --> H[初始化类]
    E --> I[解释执行初始化代码]

    F --> J[加载主类]
    F --> K[获取main方法]
    F --> L[创建参数数组]
    F --> M[创建栈帧]
    F --> N[解释执行main方法]

    N --> O[字节码解释执行循环]

    subgraph 类加载子系统
        P[类路径解析] --> Q[读取class文件]
        Q --> R[解析class文件]
        R --> S[创建运行时类结构]
        S --> T[链接:验证、准备、解析]
        T --> U[初始化]
    end

    J -.-> P

    subgraph 运行时数据区
        V[方法区] --> W[运行时常量池]
        X[JVM栈] --> Y[栈帧]
        Y --> Z[局部变量表]
        Y --> AA[操作数栈]
        AB[程序计数器]
        AC[堆]
    end

    subgraph 执行引擎
        AD[字节码解释器] --> AE[指令解码]
        AE --> AF[指令执行]
        AG[本地方法接口]
    end

    O -.-> AD
    M -.-> X
```

**主要模块说明：**

1. **启动流程**：从 `main.go` 开始，解析命令行参数，创建 JVM 实例，然后启动 JVM。

2. **类加载子系统**：负责查找、加载、解析和初始化类文件。

3. **运行时数据区**：包括方法区、堆、JVM 栈、程序计数器等内存区域。

4. **执行引擎**：包含字节码解释器和本地方法接口，负责执行字节码指令。

## 启动流程

下图展示了 JVM 启动过程中各组件之间的交互。从命令行解析到最终执行 Java 程序的 main 方法。

```mermaid
sequenceDiagram
    participant Main as main.go
    participant CMD as cmd.go
    participant JVM as jvm.go
    participant CP as classpath
    participant CL as ClassLoader
    participant Interp as interpreter.go

    Main->>CMD: parseCmd()
    CMD-->>Main: 返回命令行参数
    Main->>JVM: newJVM(cmd)
    JVM->>CP: Parse(jreOption, cpOption)
    CP-->>JVM: 返回类路径
    JVM->>CL: NewClassLoader(cp)
    CL-->>JVM: 返回类加载器
    JVM->>JVM: 创建主线程
    Main->>JVM: start()
    JVM->>JVM: initVM()
    JVM->>CL: LoadClass("sun/misc/VM")
    CL-->>JVM: 返回VM类
    JVM->>Interp: interpret(thread, verboseFlag)
    JVM->>JVM: execMain()
    JVM->>CL: LoadClass(className)
    CL-->>JVM: 返回主类
    JVM->>JVM: 获取main方法
    JVM->>JVM: 创建参数数组
    JVM->>JVM: 创建栈帧
    JVM->>Interp: interpret(thread, verboseFlag)
    Interp->>Interp: 执行字节码
```

**启动流程说明：**

1. **命令行解析**：`main.go` 调用 `cmd.go` 中的 `parseCmd()` 函数解析命令行参数。

2. **JVM 初始化**：创建 JVM 实例，解析类路径，创建类加载器和主线程。

3. **VM 初始化**：加载和初始化 `sun/misc/VM` 类，这是 Java 运行时的关键类。

4. **主类执行**：加载用户指定的主类，获取 `main` 方法，创建参数数组和栈帧，然后调用解释器执行字节码。

## 类加载流程

下图展示了类加载器的工作流程，包括类的加载、链接和初始化过程。

```mermaid
flowchart TD
    A[开始加载类] --> B{是否已加载?}
    B -->|是| C[返回已加载的类]
    B -->|否| D{是数组类?}
    D -->|是| E[loadArrayClass]
    D -->|否| F[loadNonArrayClass]

    F --> G[readClass]
    G --> H[defineClass]
    H --> I[parseClass]
    I --> J[newClass]
    H --> K[resolveSuperClass]
    H --> L[resolveInterfaces]

    E --> M[返回数组类]
    C --> N[结束]
    M --> N
    L --> O[link]
    O --> P[verify]
    O --> Q[prepare]
    O --> R[返回类]
    R --> N
```

**类加载流程说明：**

1. **检查缓存**：首先检查类是否已经加载，如果已加载则直接返回。

2. **区分类型**：判断是否是数组类，如果是则调用 `loadArrayClass`，否则调用 `loadNonArrayClass`。

3. **加载非数组类**：
   - 读取类文件数据
   - 解析类文件
   - 创建运行时类结构
   - 解析父类和接口

4. **链接**：包括验证、准备和解析三个步骤。
   - 验证：检查类文件格式是否正确
   - 准备：为静态字段分配内存并设置默认值
   - 解析：将符号引用转换为直接引用

## 字节码执行流程

下图展示了字节码解释器的执行循环，从获取当前栈帧到执行指令的完整过程。

```mermaid
flowchart TD
    A[interpret] --> B[loop]
    B --> C[获取当前栈帧]
    C --> D[获取PC]
    D --> E[设置线程PC]
    E --> F[解码指令]
    F --> G[读取操作数]
    G --> H[设置下一条指令地址]
    H --> I{是否打印日志?}
    I -->|是| J[打印指令信息]
    I -->|否| K[执行指令]
    J --> K
    K --> L{栈是否为空?}
    L -->|是| M[结束]
    L -->|否| C
```

**字节码执行流程说明：**

1. **解释器初始化**：`interpret` 函数启动解释器，进入执行循环。

2. **执行循环**：
   - 获取当前栈帧（当前正在执行的方法的栈帧）
   - 获取并设置程序计数器（PC）
   - 解码当前指令，读取操作数
   - 设置下一条指令的地址
   - 根据需要打印指令信息
   - 执行当前指令

3. **循环终止条件**：当线程的栈为空时（所有方法都执行完毕），解释器结束执行。

## 类文件解析流程

下图展示了类文件的解析过程，从魔数验证到属性读取的完整流程。

```mermaid
flowchart TD
    A[Parse] --> B[读取魔数]
    B --> C[检查魔数]
    C --> D[读取版本号]
    D --> E[检查版本号]
    E --> F[读取常量池]
    F --> G[读取访问标志]
    G --> H[读取类索引]
    H --> I[读取父类索引]
    I --> J[读取接口索引]
    J --> K[读取字段]
    K --> L[读取方法]
    L --> M[读取属性]
    M --> N[返回ClassFile]

    subgraph 常量池解析
        F1[读取常量池大小] --> F2[创建常量池数组]
        F2 --> F3[循环读取常量]
        F3 --> F4[根据tag创建常量]
        F4 --> F5[读取常量信息]
        F5 --> F6[处理特殊常量]
    end

    F -.-> F1
```

**类文件解析流程说明：**

1. **魔数验证**：首先读取并验证类文件的魔数（`0xCAFEBABE`），确保文件格式正确。

2. **版本检查**：读取并验证类文件的版本号，确保它是支持的版本。

3. **常量池解析**：读取常量池大小，创建常量池数组，然后根据每个常量的标记（tag）创建相应类型的常量对象。

4. **类信息读取**：读取类的访问标志、类索引、父类索引和接口索引。

5. **成员信息读取**：读取类的字段、方法和属性信息。

6. **返回结果**：最终返回包含完整类信息的 `ClassFile` 结构体。

## 运行时数据区结构

运行时数据区是 JVM 在运行期间用于存储数据的内存区域，包括方法区、堆、JVM 栈、本地方法栈和程序计数器。在本项目中，我们使用 Go 的数据结构来模拟这些内存区域。

```mermaid
classDiagram
    class Thread {
        +int pc
        +Stack stack
        +NewFrame(method)
        +PushFrame(frame)
        +PopFrame()
        +CurrentFrame()
    }

    class Stack {
        +uint maxSize
        +uint size
        +Frame _top
        +push(frame)
        +pop()
        +top()
    }

    class Frame {
        +Frame lower
        +LocalVars localVars
        +OperandStack operandStack
        +Thread thread
        +Method method
        +int nextPC
    }

    class LocalVars {
        +Slot[] slots
        +GetInt(index)
        +SetInt(index, val)
        +GetFloat(index)
        +SetFloat(index, val)
        +GetLong(index)
        +SetLong(index, val)
        +GetDouble(index)
        +SetDouble(index, val)
        +GetRef(index)
        +SetRef(index, ref)
    }

    class OperandStack {
        +uint size
        +Slot[] slots
        +PushInt(val)
        +PopInt()
        +PushFloat(val)
        +PopFloat()
        +PushLong(val)
        +PopLong()
        +PushDouble(val)
        +PopDouble()
        +PushRef(ref)
        +PopRef()
    }

    class Class {
        +uint16 accessFlags
        +string name
        +string superClassName
        +string[] interfaceNames
        +ConstantPool constantPool
        +Field[] fields
        +Method[] methods
        +ClassLoader loader
        +Class superClass
        +Class[] interfaces
        +uint instanceSlotCount
        +uint staticSlotCount
        +Slots staticVars
        +bool initStarted
        +Object jClass
    }

    class Method {
        +uint maxStack
        +uint maxLocals
        +byte[] code
        +ExceptionTable exceptionTable
    }

    class Object {
        +Class class
        +Object data
        +Object extra
    }

    Thread "1" --> "1" Stack : contains
    Stack "1" --> "0..*" Frame : stores
    Frame "1" --> "1" LocalVars : contains
    Frame "1" --> "1" OperandStack : contains
    Frame "1" --> "1" Thread : references
    Frame "1" --> "1" Method : executes
    Method "1" --> "1" Class : belongs_to
    Object "1" --> "1" Class : instantiates

```

**关键组件说明：**

1. **Thread（线程）**：表示 JVM 的执行线程，每个线程有自己的程序计数器（PC）和 JVM 栈。

2. **Stack（栈）**：JVM 栈，用于存储栈帧。在我们的实现中，使用链表结构来模拟栈。

3. **Frame（栈帧）**：方法执行的基本单位，包含局部变量表、操作数栈和对当前方法的引用。

4. **LocalVars（局部变量表）**：用于存储方法的参数和局部变量。

5. **OperandStack（操作数栈）**：用于存储指令操作的临时数据。

6. **Class（类）**：表示已加载的类信息，包含类的常量池、字段、方法等信息。

7. **Method（方法）**：表示类中的方法，包含方法的字节码、局部变量表大小、操作数栈大小等信息。

8. **Object（对象）**：表示堆中的对象实例，包含对象的类型和实例数据。

## 指令执行模型

JVM 的指令执行模型是基于字节码解释器的。在我们的实现中，指令执行由以下组件完成：

```mermaid
flowchart TD
    A[指令接口] --> B[无操作数指令]
    A --> C[分支指令]
    A --> D[索引8位指令]
    A --> E[索引16位指令]

    F[指令工厂] --> G[根据操作码创建指令]

    H[字节码读取器] --> I[读取uint8]
    H --> J[读取uint16]
    H --> K[读取int16]

    L[指令执行] --> M[获取操作数]
    M --> N[执行指令逻辑]
    N --> O[修改程序状态]

    subgraph 指令类型
        P[常量指令]
        Q[加载指令]
        R[存储指令]
        S[栈操作指令]
        T[数学指令]
        U[转换指令]
        V[比较指令]
        W[控制指令]
        X[引用指令]
        Y[扩展指令]
    end
```

**指令执行流程：**

1. **指令解码**：从当前方法的字节码中读取操作码（opcode）。
2. **创建指令**：根据操作码创建相应的指令对象。
3. **获取操作数**：如果指令需要操作数，从字节码中读取。
4. **执行指令**：执行指令的逻辑，可能会修改局部变量表、操作数栈或程序计数器。

**指令类型：**

我们实现了大部分 JVM 规范中定义的指令，包括常量指令、加载指令、存储指令、栈操作指令、数学指令、转换指令、比较指令、控制指令、引用指令和扩展指令。

## 本地方法接口

本地方法接口（Native Method Interface）允许 Java 代码调用非 Java 编写的方法。在我们的实现中，本地方法是用 Go 函数实现的。

```mermaid
flowchart TD
    A[本地方法注册表] --> B[注册本地方法]
    A --> C[查找本地方法]

    D[本地方法调用] --> E[查找本地方法实现]
    E --> F{找到方法?}
    F -->|是| G[调用Go函数]
    F -->|否| H[抛出UnsatisfiedLinkError]

    subgraph 标准库本地方法
        I[java.lang.Object]
        J[java.lang.Class]
        K[java.lang.String]
        L[java.lang.System]
        M[java.lang.Thread]
        N[java.io.FileDescriptor]
        O[java.io.FileInputStream]
        P[java.io.FileOutputStream]
    end
```

**本地方法实现流程：**

1. **注册本地方法**：在 JVM 启动时，我们将 Go 函数注册到本地方法注册表中。

2. **调用本地方法**：当 Java 代码调用本地方法时，执行引擎会在注册表中查找相应的 Go 函数并执行。

3. **处理结果**：本地方法执行完成后，将结果返回给 Java 代码。

我们实现了多个标准库类的本地方法，如 `Object`、`String`、`System` 等，使得 JVM 可以执行基本的 Java 程序。
