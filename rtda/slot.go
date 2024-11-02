package rtda

import "jvm-go/rtda/heap"

type Slot struct {
	num int32
	// 引用类型
	ref *heap.Object
}
