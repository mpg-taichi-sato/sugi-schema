package generator

import (
	// "encoding/json"
	// "fmt"
	// "io"
	// "io/ioutil"
	// "log"
	// "os"
	// "protoc-gen-genta/generator"
	// "protoc-gen-genta/option"
	// "strings"

	// "github.com/golang/protobuf/proto"

	"context"
	"fmt"
	"strings"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

type GenerateGoProcess struct {
}

func (p *GenerateGoProcess) Run(ctx context.Context, req *plugin.CodeGeneratorRequest) ([]*plugin.CodeGeneratorResponse_File, error) {

	g := &GoCodeGenerator{}

	files := make(map[string]*descriptor.FileDescriptorProto, len(req.ProtoFile))
	for _, f := range req.ProtoFile {
		files[f.GetName()] = f
	}
	responseFiles := make([]*plugin.CodeGeneratorResponse_File, 0, len(req.FileToGenerate))
	for _, fname := range req.FileToGenerate {
		f := files[fname]
		// get package name
		packageWords := strings.Split(f.GetPackage(), ".")
		goPackage := packageWords[len(packageWords)-1]

		protoFileInfo := CreateProtoFileInfo(f)

		messageProtos := f.GetMessageType()
		structs := make([]*GoStruct, 0, len(messageProtos))
		for i := 0; i < len(messageProtos); i++ {
			st, err := p.GetGoStruct(messageProtos[i], i, protoFileInfo, f.GetPackage())
			if err != nil {
				return nil, err
			}
			structs = append(structs, st)
		}

		goFile := &GoFile{
			PackageName:    goPackage,
			PackageComment: p.GoPackageComment(protoFileInfo),
			Structs:        structs,
		}
		content := g.Generate(ctx, goFile)
		out := strings.Replace(fname, ".proto", ".pb.go", 1)
		responseFiles = append(responseFiles, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(out),
			Content: proto.String(content),
		})
	}
	return responseFiles, nil
}

func (p *GenerateGoProcess) GoPackageComment(protoFileInfo *ProtoFileInfo) string {
	syntaxInfo := protoFileInfo.GetSyntaxInfo()
	packageInfo := protoFileInfo.GetPackageInfo()

	var sb strings.Builder
	if syntaxInfo != nil {
		sb.WriteString(syntaxInfo.GetLeadingComments())
		sb.WriteString(syntaxInfo.GetTrailingComments())
	}

	if packageInfo != nil {
		sb.WriteString(packageInfo.GetLeadingComments())
		sb.WriteString(syntaxInfo.GetTrailingComments())
	}
	return sb.String()
}

func (p *GenerateGoProcess) GetGoStruct(messageProto *descriptor.DescriptorProto, messageIndex int, protoFileInfo *ProtoFileInfo, path string) (*GoStruct, error) {
	// nested type
	nestedType := messageProto.GetNestedType()
	// mapEntry: typeName
	mapTypeNames := make(map[string]string, 0)
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
						name, err := p.GetGoTypeName(field.GetType(), field.GetTypeName())
						if err != nil {
							return nil, err
						}
						keyTypeName = name
					}

					if fieldName == "value" {
						name, err := p.GetGoTypeName(field.GetType(), field.GetTypeName())
						if err != nil {
							return nil, err
						}
						valueTypeName = name
					}
				}
				fullName := fmt.Sprintf(".%v.%v.%v", path, messageProto.GetName(), nestedMessageProto.GetName())
				mapTypeNames[fullName] = fmt.Sprintf("map[%s]%s", keyTypeName, valueTypeName)
			}

		}
	}

	messageInfo := protoFileInfo.GetMessageInfo(int32(messageIndex))
	fieldProtos := messageProto.GetField()
	fields := make([]*GoStructField, 0, len(fieldProtos))
	for i := 0; i < len(fieldProtos); i++ {
		fieldInfo := protoFileInfo.GetMessageFieldInfo(int32(messageIndex), int32(i))
		fields = append(fields, p.GetGoStructField(fieldProtos[i], fieldInfo, mapTypeNames))
	}
	st := &GoStruct{
		Name:    messageProto.GetName(),
		Fields:  fields,
		Comment: messageInfo.GetLeadingComments(),
	}

	return st, nil
}

func (p *GenerateGoProcess) GetGoStructField(fieldProto *descriptor.FieldDescriptorProto, info *descriptor.SourceCodeInfo_Location, typeNameMap map[string]string) *GoStructField {

	label := fieldProto.GetLabel()

	typ, err := p.GetGoType(fieldProto.GetType(), fieldProto.GetTypeName(), label, typeNameMap)
	if err != nil {
		panic(err)
	}

	field := &GoStructField{
		Name:    fieldProto.GetName(),
		Type:    typ,
		Comment: info.GetTrailingComments(),
	}

	return field
}

func (p *GenerateGoProcess) GetGoType(typeProto descriptor.FieldDescriptorProto_Type, typeName string, label descriptor.FieldDescriptorProto_Label, typeNameMap map[string]string) (*GoType, error) {

	if typeNameMap != nil {
		convertedTypeName, found := typeNameMap[typeName]
		if found {
			return &GoType{Name: convertedTypeName}, nil
		}
	}

	goTypeName, err := p.GetGoTypeName(typeProto, typeName)
	if err != nil {
		return nil, err
	}

	if label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
		goTypeName = "[]" + goTypeName
	}
	return &GoType{Name: goTypeName}, nil
}

func (p *GenerateGoProcess) GetGoTypeName(typeProto descriptor.FieldDescriptorProto_Type, typeName string) (string, error) {
	goTypeName := ""
	switch typeProto {
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		goTypeName = "int"
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		goTypeName = "int64"
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		goTypeName = "bool"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		goTypeName = "string"
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		goTypeName = p.GetGoStructTypeName(typeName)
	default:
		return "", fmt.Errorf("GetGoType unsupported type error %v", typeProto)
	}
	return goTypeName, nil
}

func (p *GenerateGoProcess) GetGoStructTypeName(typeName string) string {
	if typeName == ".google.protobuf.Timestamp" {
		return "time.Time"
	}
	// キャッシュで高速化できる
	typeSlice := strings.Split(typeName, ".")
	// always pointer
	return "*" + typeSlice[len(typeSlice)-1]
}
