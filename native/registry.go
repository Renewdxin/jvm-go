package native

import "jvm-go/rtda"

// NativeMethod  本地方法函数类型，接收一个帧作为参数
// 这个frame参数就是本地方法的工作空间，也就是连接Java虚拟机和Java类库的桥梁
type NativeMethod func(frame *rtda.Frame)

// registry  本地方法注册表，键为"类名~方法名~方法描述符"，值为对应的本地方法函数
var registry = map[string]NativeMethod{}

// emptyNativeMethod  空本地方法，什么也不做，用于占位
func emptyNativeMethod(frame *rtda.Frame) {
	// do nothing
}

// Register  注册本地方法
func Register(className, methodName, methodDescriptor string, method NativeMethod) {
	key := className + "~" + methodName + "~" + methodDescriptor // 构造key
	registry[key] = method                                       // 将本地方法注册到registry
}

// FindNativeMethod  查找本地方法
// 根据类名、方法名和方法描述符查找对应的本地方法函数
func FindNativeMethod(className, methodName, methodDescriptor string) NativeMethod {
	key := className + "~" + methodName + "~" + methodDescriptor
	if method, ok := registry[key]; ok { // 尝试从registry中查找
		return method
	}

	// 特殊处理：如果方法描述符为()V，且方法名为registerNatives或initIDs，则返回空本地方法
	if methodDescriptor == "()V" { // ()V表示无参数，无返回值
		if methodName == "registerNatives" || methodName == "initIDs" {
			return emptyNativeMethod
		}
	}
	return nil // 未找到，返回nil
}
