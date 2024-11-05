package main

import (
	"fmt"
	"jvm-go/classpath"
	"jvm-go/instructions/base"
	"jvm-go/rtda"
	"jvm-go/rtda/heap"
	"strings"
)

// JVM 结构体表示 Java 虚拟机。
type JVM struct {
	cmd         *Cmd              // 命令行参数
	classLoader *heap.ClassLoader // 类加载器
	mainThread  *rtda.Thread      // 主线程
}

// newJVM 创建一个新的 JVM 实例。
// cmd: 命令行参数
func newJVM(cmd *Cmd) *JVM {
	cp := classpath.Parse(cmd.XjreOption, cmd.cpOption)          // 解析类路径
	classLoader := heap.NewClassLoader(cp, cmd.verboseClassFlag) // 创建类加载器
	return &JVM{
		cmd:         cmd,
		classLoader: classLoader,
		mainThread:  rtda.NewThread(), // 创建主线程
	}
}

// start 启动 JVM。
func (vm *JVM) start() {
	vm.initVM()   // 初始化虚拟机
	vm.execMain() // 执行 main 方法
}

// initVM 初始化虚拟机。
func (vm *JVM) initVM() {
	vmClass := vm.classLoader.LoadClass("sun/misc/VM") // 加载 sun.misc.VM 类
	base.InitClass(vm.mainThread, vmClass)             // 初始化 sun.misc.VM 类，包括执行<clinit>方法
	interpret(vm.mainThread, vm.cmd.verboseInstFlag)   // 解释执行初始化类的方法，例如<clinit>
}

// execMain 执行 main 方法。
func (vm *JVM) execMain() {
	className := strings.Replace(vm.cmd.class, ".", "/", -1) // 将类名中的 "." 替换为 "/"
	mainClass := vm.classLoader.LoadClass(className)         // 加载主类
	mainMethod := mainClass.GetMainMethod()                  // 获取 main 方法
	if mainMethod == nil {                                   // 如果 main 方法不存在，则报错
		fmt.Printf("Main method not found in class %s\n", vm.cmd.class)
		return
	}

	argsArr := vm.createArgsArray()                  // 创建参数数组
	frame := vm.mainThread.NewFrame(mainMethod)      // 创建栈帧
	frame.LocalVars().SetRef(0, argsArr)             // 将参数数组存入局部变量表
	vm.mainThread.PushFrame(frame)                   // 将栈帧压入主线程的栈
	interpret(vm.mainThread, vm.cmd.verboseInstFlag) // 解释执行 main 方法
}

// createArgsArray 创建cmd参数数组。
func (vm *JVM) createArgsArray() *heap.Object {
	stringClass := vm.classLoader.LoadClass("java/lang/String") // 加载 java.lang.String 类
	argsLen := uint(len(vm.cmd.args))                           // 获取命令行参数个数
	argsArr := stringClass.ArrayClass().NewArray(argsLen)       // 创建字符串数组
	jArgs := argsArr.Refs()                                     // 获取数组元素的引用
	for i, arg := range vm.cmd.args {                           // 遍历命令行参数
		jArgs[i] = heap.JString(vm.classLoader, arg) // 将命令行参数转换为 Java 字符串，并存入数组
	}
	return argsArr
}
