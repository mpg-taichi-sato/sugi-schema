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
	"errors"
	"fmt"
	"protoc-gen-genta/option"
	"strings"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

type GenerateGoProcess struct {
	TagJSON bool
}

func (p *GenerateGoProcess) Run(ctx context.Context, req *plugin.CodeGeneratorRequest) ([]*plugin.CodeGeneratorResponse_File, error) {

	g := &GoCodeGenerator{}

	protoTypeInfo := CreateProtoTypeInfo(&GoTypeConverter{})

	files := make(map[string]*descriptor.FileDescriptorProto, len(req.ProtoFile))
	for i := 0; i < len(req.ProtoFile); i++ {
		f := req.ProtoFile[i]
		files[f.GetName()] = req.ProtoFile[i]
		err := protoTypeInfo.UpdateTypeMapByFileProto(f, i)
		if err != nil {
			return nil, err
		}
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
			st, err := p.GetGoStruct(messageProtos[i], i, protoFileInfo, protoTypeInfo)
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

func (p *GenerateGoProcess) GetGoStruct(messageProto *descriptor.DescriptorProto, messageIndex int, protoFileInfo *ProtoFileInfo, protoTypeInfo *ProtoTypeInfo) (*GoStruct, error) {

	messageInfo := protoFileInfo.GetMessageInfo(int32(messageIndex))
	fieldProtos := messageProto.GetField()
	fields := make([]*GoStructField, 0, len(fieldProtos))
	for i := 0; i < len(fieldProtos); i++ {
		fieldProto := fieldProtos[i]
		fieldInfo := protoFileInfo.GetMessageFieldInfo(int32(messageIndex), int32(i))
		// fields = append(fields, p.GetGoStructField(fieldProtos[i], fieldInfo, mapTypeNames))

		dataType, err := protoTypeInfo.GetTypeName(fieldProto)
		if err != nil {
			return nil, err
		}

		tagOption, found := option.GetGoTagOption(fieldProto.GetOptions())
		if !found && p.TagJSON {
			tagOption = fmt.Sprintf("json:\"%s\"", fieldProto.GetName())
		}

		fields = append(fields, &GoStructField{
			Name:    fieldProto.GetName(),
			Type:    &GoType{Name: dataType},
			Comment: fieldInfo.GetTrailingComments(),
			Tag:     tagOption,
		})
	}
	st := &GoStruct{
		Name:    messageProto.GetName(),
		Fields:  fields,
		Comment: messageInfo.GetLeadingComments(),
	}

	return st, nil
}

type GoTypeConverter struct {
}

func (c *GoTypeConverter) FormatMapType(key, value string) string {
	return fmt.Sprintf("map[%s]%s", key, value)
}

func (c *GoTypeConverter) FormatArrayType(srcType string) string {
	return "[]" + srcType
}

func (c *GoTypeConverter) GetScalarTypeName(typeProto descriptor.FieldDescriptorProto_Type) (string, error) {
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
		return "", errors.New("not scalar")
	default:
		return "", fmt.Errorf("GetGoType unsupported type error %v", typeProto)
	}
	return goTypeName, nil
}

func (c *GoTypeConverter) GetMessageTypeName(protoTypeName string) string {
	if protoTypeName == ".google.protobuf.Timestamp" {
		return "time.Time"
	}

	typeSlice := strings.Split(protoTypeName, ".")
	// always pointer
	return "*" + typeSlice[len(typeSlice)-1]
}
