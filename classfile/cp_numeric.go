package classfile

import "math"

/*
	CONSTANT_Integer_info {
	    u1 tag;
	    u4 bytes;
	}
*/
type ConstantIntegerInfo struct {
	val int32
}

func (cii *ConstantIntegerInfo) readInfo(reader *ClassReader) {
	bytes := reader.readUint32()
	cii.val = int32(bytes)
}
func (cii *ConstantIntegerInfo) Value() int32 {
	return cii.val
}

/*
	CONSTANT_Float_info {
	    u1 tag;
	    u4 bytes;
	}
*/
type ConstantFloatInfo struct {
	val float32
}

func (cii *ConstantFloatInfo) readInfo(reader *ClassReader) {
	bytes := reader.readUint32()
	cii.val = math.Float32frombits(bytes)
}
func (cii *ConstantFloatInfo) Value() float32 {
	return cii.val
}

/*
	CONSTANT_Long_info {
	    u1 tag;
	    u4 high_bytes;
	    u4 low_bytes;
	}
*/
type ConstantLongInfo struct {
	val int64
}

func (cii *ConstantLongInfo) readInfo(reader *ClassReader) {
	bytes := reader.readUint64()
	cii.val = int64(bytes)
}
func (cii *ConstantLongInfo) Value() int64 {
	return cii.val
}

/*
	CONSTANT_Double_info {
	    u1 tag;
	    u4 high_bytes;
	    u4 low_bytes;
	}
*/
type ConstantDoubleInfo struct {
	val float64
}

func (cii *ConstantDoubleInfo) readInfo(reader *ClassReader) {
	bytes := reader.readUint64()
	cii.val = math.Float64frombits(bytes)
}
func (cii *ConstantDoubleInfo) Value() float64 {
	return cii.val
}
