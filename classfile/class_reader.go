package classfile

import "encoding/binary"

// ClassReader 用于读取类文件的数据。
// 不同于使用索引记录数据位置的方式，ClassReader 利用 Go 语言的切片特性（reslice）
// 来跳过已读取的数据，高效地管理数据读取过程。

type ClassReader struct {
	data []byte // 存储类文件的字节数据
}

// readUint8 读取一个 uint8 类型的数据。
func (cr *ClassReader) readUint8() uint8 {
	val := cr.data[0]     // 读取第一个字节
	cr.data = cr.data[1:] // 将切片起始位置后移一个字节，相当于跳过已读取的数据
	return val
}

// readUint16 读取一个 uint16 (u2) 类型的数据。
func (cr *ClassReader) readUint16() uint16 {
	val := binary.BigEndian.Uint16(cr.data) // 使用 BigEndian 字节序解码数据
	cr.data = cr.data[2:]                   // 后移两个字节
	return val
}

// readUint32 读取一个 uint32 (u4) 类型的数据。
func (cr *ClassReader) readUint32() uint32 {
	val := binary.BigEndian.Uint32(cr.data) // 使用 BigEndian 字节序解码数据
	cr.data = cr.data[4:]                   // 后移四个字节
	return val
}

// readUint64 读取一个 uint64 类型的数据。
func (cr *ClassReader) readUint64() uint64 {
	val := binary.BigEndian.Uint64(cr.data) // 使用 BigEndian 字节序解码数据
	cr.data = cr.data[8:]                   // 后移八个字节
	return val
}

// readUint16s 读取一个 uint16 数组。
// 首先读取数组长度 (u2)，然后读取指定数量的 uint16 元素。
func (cr *ClassReader) readUint16s() []uint16 {
	n := cr.readUint16()   // 读取数组长度
	s := make([]uint16, n) // 创建指定长度的数组
	for i := range s {
		s[i] = cr.readUint16() // 读取每个元素
	}
	return s
}

// readBytes 读取指定长度的字节数组。
//
// 参数：
//
//	n: 要读取的字节数。
//
// 返回值：
//
//	[]byte: 读取到的字节数组。
func (cr *ClassReader) readBytes(n uint32) []byte {
	bytes := cr.data[:n]  // 获取指定长度的字节切片
	cr.data = cr.data[n:] // 后移 n 个字节
	return bytes
}
