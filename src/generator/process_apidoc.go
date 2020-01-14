package generator

import (
	"context"
	"errors"
	"protoc-gen-genta/option"
	"strings"

	"github.com/golang/protobuf/proto"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

type GenerateAPIDocProcess struct {
}

func (p *GenerateAPIDocProcess) Run(ctx context.Context, req *plugin.CodeGeneratorRequest) ([]*plugin.CodeGeneratorResponse_File, error) {

	g := &APIDocGenerator{}
	fileProtos := make(map[string]*descriptor.FileDescriptorProto, len(req.ProtoFile))
	protoFileInfoList := make([]*ProtoFileInfo, 0, len(req.ProtoFile))

	protoTypeInfo := CreateProtoTypeInfo(&GoTypeConverter{})

	for i := 0; i < len(req.ProtoFile); i++ {
		f := req.ProtoFile[i]
		fileProtos[f.GetName()] = req.ProtoFile[i]
		err := protoTypeInfo.UpdateTypeMapByFileProto(f, i)
		if err != nil {
			return nil, err
		}
		protoFileInfoList = append(protoFileInfoList, CreateProtoFileInfo(f))
	}

	responseFiles := make([]*plugin.CodeGeneratorResponse_File, 0)

	for _, fname := range req.FileToGenerate {
		// apiを含むファイルのみ出力
		f := fileProtos[fname]
		protoFileInfo := CreateProtoFileInfo(f)
		serviceProtos := f.GetService()
		services := make([]*APIDocService, 0)
		for i := 0; i < len(serviceProtos); i++ {
			fileServices, err := p.GetService(serviceProtos[i], i, protoFileInfo, protoFileInfoList, protoTypeInfo)
			if err != nil {
				return nil, err
			}
			services = append(services, fileServices...)
		}

		if len(services) == 0 {
			continue
		}

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

		apiDocInfo := &APIDocInfo{
			HeadComment: sb.String(),
			Services:    services,
		}

		content := g.Generate(ctx, apiDocInfo)
		out := strings.Replace(fname, ".proto", ".md", 1)
		responseFiles = append(responseFiles, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(out),
			Content: proto.String(content),
		})
	}

	return responseFiles, nil
}

func (p *GenerateAPIDocProcess) GetService(serviceProto *descriptor.ServiceDescriptorProto, serviceIndex int, protoFileInfo *ProtoFileInfo, protoFileInfoList []*ProtoFileInfo, protoTypeInfo *ProtoTypeInfo) ([]*APIDocService, error) {
	// serviceInfo := protoFileInfo.GetServiceInfo(int32(serviceIndex))
	methodProtos := serviceProto.GetMethod()
	apiList := make([]*APIDocService, 0, len(methodProtos))
	for i := 0; i < len(methodProtos); i++ {
		methodProto := methodProtos[i]
		apiOption := option.GetAPIOption(methodProto.GetOptions())
		methodInfo := protoFileInfo.GetMethodInfo(int32(serviceIndex), int32(i))

		api := &APIDocService{
			Method:  apiOption.GetMethod(),
			Path:    apiOption.GetPath(),
			Comment: methodInfo.GetLeadingComments(),
		}

		if methodProto.InputType != nil && methodProto.GetInputType() != ".google.protobuf.Empty" {
			protoMessage, ok := protoTypeInfo.ProtoMessage(*methodProto.InputType)
			if !ok {
				return nil, errors.New("message not found")
			}
			model, err := p.GetAPIDocStructModel(protoMessage.DescriptorProto, protoMessage.Index, protoFileInfoList[protoMessage.FileIndex], protoTypeInfo)
			if err != nil {
				return nil, err
			}
			api.Request = model
		}

		if methodProto.OutputType != nil {
			protoMessage, ok := protoTypeInfo.ProtoMessage(*methodProto.OutputType)
			if !ok {
				return nil, errors.New("message not found")
			}
			model, err := p.GetAPIDocStructModel(protoMessage.DescriptorProto, protoMessage.Index, protoFileInfoList[protoMessage.FileIndex], protoTypeInfo)
			if err != nil {
				return nil, err
			}
			api.Response = model
		}
		apiList = append(apiList, api)
	}
	return apiList, nil
}

func (p *GenerateAPIDocProcess) GetAPIDocStructModel(messageProto *descriptor.DescriptorProto, messageIndex int, protoFileInfo *ProtoFileInfo, protoTypeInfo *ProtoTypeInfo) (*APIDocStructModel, error) {

	messageInfo := protoFileInfo.GetMessageInfo(int32(messageIndex))
	fieldProtos := messageProto.GetField()
	fields := make([]*APIDocStructField, 0, len(fieldProtos))
	for i := 0; i < len(fieldProtos); i++ {
		fieldProto := fieldProtos[i]
		fieldInfo := protoFileInfo.GetMessageFieldInfo(int32(messageIndex), int32(i))
		dataType, err := protoTypeInfo.GetTypeName(fieldProto)
		if err != nil {
			return nil, err
		}

		fields = append(fields, &APIDocStructField{
			Name:        fieldProto.GetName(),
			Description: strings.ReplaceAll(fieldInfo.GetTrailingComments(), "\n", " "),
			DataType:    dataType,
		})
	}
	st := &APIDocStructModel{
		Name:    messageProto.GetName(),
		Fields:  fields,
		Comment: messageInfo.GetLeadingComments(),
	}

	return st, nil
}
