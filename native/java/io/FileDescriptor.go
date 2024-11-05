package io

import (
	"jvm-go/native"
	"jvm-go/rtda"
	"os"
	"syscall"
)

const fd = "java/io/FileDescriptor"

func init() {
	native.Register(fd, "set", "(I)J", set)
}

// private static native long set(int d);
// (I)J
func set(frame *rtda.Frame) {
	fileDescriptor := frame.LocalVars().GetInt(0) // 获取第一个局部变量，即文件描述符

	// 将 Go 的文件描述符转换为 long 并推入操作数栈
	switch fileDescriptor {
	case 0: // 标准输入
		frame.OperandStack().PushLong(int64(syscall.Stdin))
	case 1: // 标准输出
		frame.OperandStack().PushLong(int64(syscall.Stdout))
	case 2: // 标准错误
		frame.OperandStack().PushLong(int64(syscall.Stderr))
	default:
		// 处理其他文件描述符的情况，可能需要查找已打开的文件或抛出异常
		// todo 一个简单的示例，创建一个新的文件描述符，实际应用中需要根据具体情况处理
		f, err := os.CreateTemp("", "jvmgo-fd")
		if err != nil {
			panic(err) // 处理错误，例如抛出 Java 异常
		}
		defer f.Close()
		frame.OperandStack().PushLong(int64(f.Fd()))
	}
}
