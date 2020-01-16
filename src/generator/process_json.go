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
	"encoding/json"
	"protoc-gen-genta/option"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

type GenerateJSONProcess struct {
}

func (p *GenerateJSONProcess) Run(ctx context.Context, req *plugin.CodeGeneratorRequest) ([]*plugin.CodeGeneratorResponse_File, error) {

	fileProtos := make(map[string]*descriptor.FileDescriptorProto, len(req.ProtoFile))
	for i, f := range req.ProtoFile {
		fileProtos[f.GetName()] = req.ProtoFile[i]
	}

	responseFiles := make([]*plugin.CodeGeneratorResponse_File, 0, len(req.FileToGenerate))

	for _, fname := range req.FileToGenerate {

		f := fileProtos[fname]
		protoFileInfo := CreateProtoFileInfo(f)
		messageProtos := f.GetMessageType()
		messages := make([]map[string]interface{}, 0, len(messageProtos))
		for i := 0; i < len(messageProtos); i++ {
			messages = append(messages, p.GetMessage(messageProtos[i], i, protoFileInfo))
		}

		serviceProtos := f.GetService()
		services := make([]map[string]interface{}, 0, len(serviceProtos))
		for i := 0; i < len(serviceProtos); i++ {
			services = append(services, p.GetService(serviceProtos[i], i, protoFileInfo))
		}

		dataMap := map[string]interface{}{
			"filename": fname,
			"protofile": map[string]interface{}{
				"name":       f.GetName(),
				"package":    f.GetPackage(),
				"Dependency": f.GetDependency(),
				"messages":   messages,
				"services":   services,
			},
			// "file": f,
			// "sourcecodeinfolocations": protoFileInfo.sourceCodeInfoLocations,
		}
		syntaxInfo := protoFileInfo.GetSyntaxInfo()
		if syntaxInfo != nil {
			dataMap["syntaxLeadingComments"] = syntaxInfo.GetLeadingComments()
			dataMap["syntaxTrailingComments"] = syntaxInfo.GetTrailingComments()
		}

		packageInfo := protoFileInfo.GetPackageInfo()
		if packageInfo != nil {
			dataMap["packageLeadingComments"] = packageInfo.GetLeadingComments()
			dataMap["packageTrailingComments"] = packageInfo.GetTrailingComments()
		}

		i, _ := json.MarshalIndent(dataMap, "", "   ")
		dataJSON := string(i)
		out := fname + ".json"
		responseFiles = append(responseFiles, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(out),
			Content: proto.String(dataJSON),
		})
	}
	return responseFiles, nil
}

func (p *GenerateJSONProcess) GetService(serviceProto *descriptor.ServiceDescriptorProto, serviceIndex int, protoFileInfo *ProtoFileInfo) map[string]interface{} {

	serviceInfo := protoFileInfo.GetServiceInfo(int32(serviceIndex))

	methodProtos := serviceProto.GetMethod()
	methods := make([]map[string]interface{}, 0, len(methodProtos))
	for i := 0; i < len(methodProtos); i++ {
		methodProto := methodProtos[i]
		methodOptions := methodProto.GetOptions()
		apiOption := option.GetAPIOption(methodOptions)
		methodInfo := protoFileInfo.GetMethodInfo(int32(serviceIndex), int32(i))
		methods = append(methods, map[string]interface{}{
			"name":                    methodProto.GetName(),
			"InputType":               methodProto.GetInputType(),
			"OutputType":              methodProto.GetOutputType(),
			"method":                  apiOption.GetMethod(),
			"path":                    apiOption.GetPath(),
			"leadingComments":         methodInfo.GetLeadingComments(),
			"trailingComments":        methodInfo.GetTrailingComments(),
			"leadingDetachedComments": methodInfo.GetLeadingDetachedComments(),
		})
	}
	service := map[string]interface{}{
		"name":                    serviceProto.GetName(),
		"method":                  methods,
		"leadingComments":         serviceInfo.GetLeadingComments(),
		"trailingComments":        serviceInfo.GetTrailingComments(),
		"leadingDetachedComments": serviceInfo.GetLeadingDetachedComments(),
	}
	return service
}

func (p *GenerateJSONProcess) GetMessage(messageProto *descriptor.DescriptorProto, messageIndex int, protoFileInfo *ProtoFileInfo) map[string]interface{} {

	messageInfo := protoFileInfo.GetMessageInfo(int32(messageIndex))

	fieldProtos := messageProto.GetField()
	fields := make([]map[string]interface{}, 0, len(fieldProtos))
	for i := 0; i < len(fieldProtos); i++ {
		fieldInfo := protoFileInfo.GetMessageFieldInfo(int32(messageIndex), int32(i))
		fields = append(fields, p.GetField(fieldProtos[i], fieldInfo))
	}
	message := map[string]interface{}{
		"name":                    messageProto.GetName(),
		"leadingComments":         messageInfo.GetLeadingComments(),
		"trailingComments":        messageInfo.GetTrailingComments(),
		"leadingDetachedComments": messageInfo.GetLeadingDetachedComments(),
		"nested_type":             messageProto.GetNestedType(),
	}
	if len(fields) != 0 {
		message["fields"] = fields
	}

	return message
}

func (p *GenerateJSONProcess) GetField(fieldProto *descriptor.FieldDescriptorProto, info *descriptor.SourceCodeInfo_Location) map[string]interface{} {
	label := fieldProto.GetLabel()
	tagOption, _ := option.GetGoTagOption(fieldProto.GetOptions())
	field := map[string]interface{}{
		"name":                    fieldProto.GetName(),
		"type":                    fieldProto.GetType(), // 11ならmessage
		"typename":                fieldProto.GetTypeName(),
		"repeated":                label == descriptor.FieldDescriptorProto_LABEL_REPEATED,
		"leadingComments":         info.GetLeadingComments(),
		"trailingComments":        info.GetTrailingComments(),
		"leadingDetachedComments": info.GetLeadingDetachedComments(),
		"go_tag":                  tagOption,
	}

	return field
}
