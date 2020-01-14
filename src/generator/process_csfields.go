package generator

import (
	"context"
	"strings"

	"github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

// GenerateCSFieldsProcess comma separated fields
type GenerateCSFieldsProcess struct {
}

func (p *GenerateCSFieldsProcess) Run(ctx context.Context, req *plugin.CodeGeneratorRequest) ([]*plugin.CodeGeneratorResponse_File, error) {
	fileProtos := make(map[string]*descriptor.FileDescriptorProto, len(req.ProtoFile))
	for i, f := range req.ProtoFile {
		fileProtos[f.GetName()] = req.ProtoFile[i]
	}

	responseFiles := make([]*plugin.CodeGeneratorResponse_File, 0, len(req.FileToGenerate))
	for _, fname := range req.FileToGenerate {
		f := fileProtos[fname]

		messageProtos := f.GetMessageType()
		var sb strings.Builder
		for i := 0; i < len(messageProtos); i++ {
			messageProto := messageProtos[i]
			fieldProtos := messageProto.GetField()
			sb.WriteString("Message_" + messageProto.GetName() + "\n")
			fieldNames := make([]string, 0, len(fieldProtos))
			for j := 0; j < len(fieldProtos); j++ {
				fieldNames = append(fieldNames, fieldProtos[j].GetName())
			}
			sb.WriteString(strings.Join(fieldNames, ","))
			sb.WriteString("\n")
		}

		out := fname + ".txt"
		responseFiles = append(responseFiles, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(out),
			Content: proto.String(sb.String()),
		})
	}
	return responseFiles, nil
}
