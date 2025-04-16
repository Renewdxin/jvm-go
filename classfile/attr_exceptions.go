package classfile

/*
	Exceptions_attribute {
	    u2 attribute_name_index;
	    u4 attribute_length;
	    u2 number_of_exceptions;
	    u2 exception_index_table[number_of_exceptions];
	}
*/
type ExceptionsAttribute struct {
	exceptionIndexTable []uint16
}

func (e *ExceptionsAttribute) readInfo(reader *ClassReader) {
	e.exceptionIndexTable = reader.readUint16s()
}

func (e *ExceptionsAttribute) ExceptionIndexTable() []uint16 {
	return e.exceptionIndexTable
}
