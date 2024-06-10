package heap

import "strings"

type MethodDescriptorParser struct {
	raw    string
	offset int
	parsed *MethodDescriptor
}

func parseMethodDescriptor(descriptor string) *MethodDescriptor {
	parser := &MethodDescriptorParser{}
	return parser.parse(descriptor)
}

func (mdp *MethodDescriptorParser) parse(descriptor string) *MethodDescriptor {
	mdp.raw = descriptor
	mdp.parsed = &MethodDescriptor{}
	mdp.startParams()
	mdp.parseParamTypes()
	mdp.endParams()
	mdp.parseReturnType()
	mdp.finish()
	return mdp.parsed
}

func (mdp *MethodDescriptorParser) startParams() {
	if mdp.readUint8() != '(' {
		mdp.causePanic()
	}
}
func (mdp *MethodDescriptorParser) endParams() {
	if mdp.readUint8() != ')' {
		mdp.causePanic()
	}
}
func (mdp *MethodDescriptorParser) finish() {
	if mdp.offset != len(mdp.raw) {
		mdp.causePanic()
	}
}

func (mdp *MethodDescriptorParser) causePanic() {
	panic("BAD descriptor: " + mdp.raw)
}

func (mdp *MethodDescriptorParser) readUint8() uint8 {
	b := mdp.raw[mdp.offset]
	mdp.offset++
	return b
}
func (mdp *MethodDescriptorParser) unreadUint8() {
	mdp.offset--
}

func (mdp *MethodDescriptorParser) parseParamTypes() {
	for {
		t := mdp.parseFieldType()
		if t != "" {
			mdp.parsed.addParameterType(t)
		} else {
			break
		}
	}
}

func (mdp *MethodDescriptorParser) parseReturnType() {
	if mdp.readUint8() == 'V' {
		mdp.parsed.returnType = "V"
		return
	}

	mdp.unreadUint8()
	t := mdp.parseFieldType()
	if t != "" {
		mdp.parsed.returnType = t
		return
	}

	mdp.causePanic()
}

func (mdp *MethodDescriptorParser) parseFieldType() string {
	switch mdp.readUint8() {
	case 'B':
		return "B"
	case 'C':
		return "C"
	case 'D':
		return "D"
	case 'F':
		return "F"
	case 'I':
		return "I"
	case 'J':
		return "J"
	case 'S':
		return "S"
	case 'Z':
		return "Z"
	case 'L':
		return mdp.parseObjectType()
	case '[':
		return mdp.parseArrayType()
	default:
		mdp.unreadUint8()
		return ""
	}
}

func (mdp *MethodDescriptorParser) parseObjectType() string {
	unread := mdp.raw[mdp.offset:]
	semicolonIndex := strings.IndexRune(unread, ';')
	if semicolonIndex == -1 {
		mdp.causePanic()
		return ""
	} else {
		objStart := mdp.offset - 1
		objEnd := mdp.offset + semicolonIndex + 1
		mdp.offset = objEnd
		descriptor := mdp.raw[objStart:objEnd]
		return descriptor
	}
}

func (mdp *MethodDescriptorParser) parseArrayType() string {
	arrStart := mdp.offset - 1
	mdp.parseFieldType()
	arrEnd := mdp.offset
	descriptor := mdp.raw[arrStart:arrEnd]
	return descriptor
}
