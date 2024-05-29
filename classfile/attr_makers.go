package classfile

type DeprecatedAttribute struct { MarkerAttribute }
type SyntheticAttribute struct { MarkerAttribute }

type MarkerAttribute struct{}

// have no data
func (self *MarkerAttribute) readInfo(reader *ClassReader) {
// read nothing
}