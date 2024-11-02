package classfile

import "fmt"

// ConstantPool 常量池，存储类文件中各种常量信息，例如字符串字面量、类名、方法名等。
// ConstantPool 是一个切片，存储 ConstantInfo 接口的实现。
type ConstantPool []ConstantInfo

// readConstantPool 从 ClassReader 中读取常量池信息。
//
// 参数：
//
//	reader: ClassReader，用于读取类文件数据。
//
// 返回值：
//
//	ConstantPool: 读取到的常量池。
func readConstantPool(reader *ClassReader) ConstantPool {
	cpCount := int(reader.readUint16()) // 读取常量池大小
	cp := make([]ConstantInfo, cpCount) // 创建 ConstantPool 切片

	// 常量池索引从 1 开始，到 constant_pool_count - 1 结束。索引 0 保留。
	for i := 1; i < cpCount; i++ {
		cp[i] = readConstantInfo(reader, cp) // 读取单个常量信息

		// 处理CONSTANT_Long_info 和 CONSTANT_Double_info 的特殊情况
		// http://docs.oracle.com/javase/specs/jvms/se8/html/jvms-4.html#jvms-4.4.5
		// 8 字节常量（long 和 double）占用常量池中的两个条目。
		// 如果索引 n 的常量是 CONSTANT_Long_info 或 CONSTANT_Double_info，
		// 那么下一个可用常量的索引是 n+2。索引 n+1 虽然有效，但不可用。
		switch cp[i].(type) {
		case *ConstantLongInfo, *ConstantDoubleInfo:
			i++ // 跳过下一个索引
		}
	}

	return cp
}

// getConstantInfo 获取指定索引的常量信息。
//
// 参数：
//
//	index: 常量池索引。
//
// 返回值：
//
//	ConstantInfo: 指定索引的常量信息。
//
// 异常：
//
//	如果索引无效，则抛出 panic。
func (cp ConstantPool) getConstantInfo(index uint16) ConstantInfo {
	if cpInfo := cp[index]; cpInfo != nil { // 检查索引是否有效
		return cpInfo
	}
	panic(fmt.Errorf("无效的常量池索引: %v!", index)) // 抛出错误信息
}

// getNameAndType 获取指定索引的名称和类型描述符。
// 该索引指向一个 CONSTANT_NameAndType_info 常量。
//
// 参数：
//
//	index: 常量池索引。
//
// 返回值：
//
//	string, string: 名称和类型描述符。
func (cp ConstantPool) getNameAndType(index uint16) (string, string) {
	ntInfo := cp.getConstantInfo(index).(*ConstantNameAndTypeInfo) // 获取 CONSTANT_NameAndType_info 常量
	name := cp.getUtf8(ntInfo.nameIndex)                           // 获取名称
	_type := cp.getUtf8(ntInfo.descriptorIndex)                    // 获取类型描述符
	return name, _type
}

// getClassName 获取指定索引的类名。
// 该索引指向一个 CONSTANT_Class_info 常量。
//
// 参数：
//
//	index: 常量池索引。
//
// 返回值：
//
//	string: 类名。
func (cp ConstantPool) getClassName(index uint16) string {
	classInfo := cp.getConstantInfo(index).(*ConstantClassInfo) // 获取 CONSTANT_Class_info 常量
	return cp.getUtf8(classInfo.nameIndex)                      // 获取类名
}

// getUtf8 获取指定索引的 UTF-8 字符串。
// 该索引指向一个 CONSTANT_Utf8_info 常量。
//
// 参数：
//
//	index: 常量池索引。
//
// 返回值：
//
//	string: UTF-8 字符串。
func (cp ConstantPool) getUtf8(index uint16) string {
	utf8Info := cp.getConstantInfo(index).(*ConstantUtf8Info) // 获取 CONSTANT_Utf8_info 常量
	return utf8Info.str                                       // 返回字符串
}
