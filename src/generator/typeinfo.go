package generator

import (
	"fmt"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// ProtoTypeInfo mapがnestedTypeに含まれているということを意識せず使えるようにする
type ProtoTypeInfo struct {
	TypeConverter TypeConverter
	TypeMap       map[string]*ProtoMessage
	MapEntries    map[string]string // mapEntry: typeName
}

func CreateProtoTypeInfo(typeConverter TypeConverter) *ProtoTypeInfo {
	return &ProtoTypeInfo{
		TypeConverter: typeConverter,
		TypeMap:       make(map[string]*ProtoMessage, 0),
		MapEntries:    make(map[string]string, 0),
	}
}

type ProtoMessage struct {
	FileIndex       int
	Index           int
	SingleTypeName  string
	DescriptorProto *descriptor.DescriptorProto
}

type TypeConverter interface {
	FormatMapType(key, value string) string
	FormatArrayType(srcType string) string
	GetScalarTypeName(typeProto descriptor.FieldDescriptorProto_Type) (string, error)
	GetMessageTypeName(protoTypeName string) string
}

// getSingleTypeName not array or map
func (info *ProtoTypeInfo) getSingleTypeName(fieldProto *descriptor.FieldDescriptorProto) (string, error) {
	typeProto := fieldProto.GetType()
	if typeProto == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		return info.TypeConverter.GetMessageTypeName(fieldProto.GetTypeName()), nil
	}

	return info.TypeConverter.GetScalarTypeName(typeProto)
}

func (info *ProtoTypeInfo) UpdateTypeMapByFileProto(fileProto *descriptor.FileDescriptorProto, fileIndex int) error {
	packageName := fileProto.GetPackage()
	messageProtos := fileProto.GetMessageType()
	for i, messageProto := range messageProtos {
		// self
		selfPath := fmt.Sprintf(".%s.%s", packageName, messageProto.GetName())
		info.TypeMap[selfPath] = &ProtoMessage{
			FileIndex:       fileIndex,
			Index:           i,
			SingleTypeName:  info.TypeConverter.GetMessageTypeName(selfPath),
			DescriptorProto: messageProtos[i],
		}

		// nested
		nestedType := messageProto.GetNestedType()

		if nestedType != nil {
			for _, nestedMessageProto := range nestedType {
				nestedMessageOptions := nestedMessageProto.GetOptions()
				if nestedMessageOptions == nil {
					continue
				}
				// is this message map
				if nestedMessageOptions.GetMapEntry() {
					mapEntryFieldProtos := nestedMessageProto.GetField()
					keyTypeName := ""
					valueTypeName := ""
					for _, field := range mapEntryFieldProtos {
						fieldName := field.GetName()
						if fieldName == "key" {
							name, err := info.getSingleTypeName(field)
							if err != nil {
								return err
							}
							keyTypeName = name
						}

						if fieldName == "value" {
							name, err := info.getSingleTypeName(field)
							if err != nil {
								return err
							}
							valueTypeName = name
						}
					}
					path := fmt.Sprintf(".%v.%v.%v", packageName, messageProto.GetName(), nestedMessageProto.GetName())
					info.MapEntries[path] = info.TypeConverter.FormatMapType(keyTypeName, valueTypeName)
				}

			}
		}
	}
	return nil
}

func (info *ProtoTypeInfo) ProtoMessage(protoTypePath string) (*ProtoMessage, bool) {
	rtn, ok := info.TypeMap[protoTypePath]
	return rtn, ok
}

func (info *ProtoTypeInfo) GetTypeName(fieldProto *descriptor.FieldDescriptorProto) (string, error) {
	protoTypeName := fieldProto.GetTypeName()
	label := fieldProto.GetLabel()

	isArray := false
	if label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
		// map
		typeName, found := info.MapEntries[protoTypeName]
		if found {
			return typeName, nil
		}
		isArray = true
	}

	typeName := ""
	protoMessage, found := info.TypeMap[protoTypeName]
	if !found {
		tn, err := info.getSingleTypeName(fieldProto)
		if err != nil {
			return "", err
		}
		typeName = tn
	} else {
		typeName = protoMessage.SingleTypeName
	}
	if typeName == "Timestamp" {
		panic(fmt.Errorf("Timestamp %v %v", fieldProto.GetTypeName(), found))
	}

	if isArray {
		typeName = info.TypeConverter.FormatArrayType(typeName)
	}
	return typeName, nil
}
