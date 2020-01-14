package generator

import (
	"testing"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/google/go-cmp/cmp"
)

func TestGenerateGoProcess_GetGoStruct(t *testing.T) {
	type args struct {
		messageProto  *descriptor.DescriptorProto
		messageIndex  int
		protoFileInfo *ProtoFileInfo
		protoTypeInfo *ProtoTypeInfo
	}

	strPointer := func(v string) *string {
		return &v
	}
	// int32Pointer := func(v int32) *int32 {
	// 	return &v
	// }

	boolPointer := func(v bool) *bool {
		return &v
	}

	fieldTypePointer := func(v descriptor.FieldDescriptorProto_Type) *descriptor.FieldDescriptorProto_Type {
		return &v
	}

	fieldLabelPointer := func(v descriptor.FieldDescriptorProto_Label) *descriptor.FieldDescriptorProto_Label {
		return &v
	}

	protoTypeInfo := CreateProtoTypeInfo(&GoTypeConverter{})
	protoTypeInfo.MapEntries[".package.DummyMessage.DummyMapEntry"] = "map[int]string"

	tests := []struct {
		name    string
		p       *GenerateGoProcess
		args    args
		want    *GoStruct
		wantErr bool
	}{
		{
			name: "正常系 map",
			p:    &GenerateGoProcess{},
			args: args{
				messageProto: &descriptor.DescriptorProto{
					Name: strPointer("DummyMessage"),
					Field: []*descriptor.FieldDescriptorProto{
						&descriptor.FieldDescriptorProto{
							Name:     strPointer("dummyMap"),
							Label:    fieldLabelPointer(descriptor.FieldDescriptorProto_LABEL_REPEATED),
							Type:     fieldTypePointer(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
							TypeName: strPointer(".package.DummyMessage.DummyMapEntry"),
						},
					},
					NestedType: []*descriptor.DescriptorProto{
						&descriptor.DescriptorProto{
							Name: strPointer("DummyMapEntry"),
							Field: []*descriptor.FieldDescriptorProto{
								&descriptor.FieldDescriptorProto{
									Name: strPointer("key"),
									Type: fieldTypePointer(descriptor.FieldDescriptorProto_TYPE_INT32),
								},
								&descriptor.FieldDescriptorProto{
									Name: strPointer("value"),
									Type: fieldTypePointer(descriptor.FieldDescriptorProto_TYPE_STRING),
								},
							},
							Options: &descriptor.MessageOptions{
								MapEntry: boolPointer(true),
							},
						},
					},
				},
				messageIndex:  1,
				protoFileInfo: &ProtoFileInfo{},
				protoTypeInfo: protoTypeInfo,
			},
			want: &GoStruct{
				Name: "DummyMessage",
				Fields: []*GoStructField{
					&GoStructField{
						Name: "dummyMap",
						Type: &GoType{
							Name: "map[int]string",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "正常系 timestamp",
			p:    &GenerateGoProcess{},
			args: args{
				messageProto: &descriptor.DescriptorProto{
					Name: strPointer("DummyMessage"),
					Field: []*descriptor.FieldDescriptorProto{
						&descriptor.FieldDescriptorProto{
							Name:     strPointer("CreatedAt"),
							Type:     fieldTypePointer(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
							TypeName: strPointer(".google.protobuf.Timestamp"),
						},
					},
				},
				messageIndex:  1,
				protoFileInfo: &ProtoFileInfo{},
				protoTypeInfo: protoTypeInfo,
			},
			want: &GoStruct{
				Name: "DummyMessage",
				Fields: []*GoStructField{
					&GoStructField{
						Name: "CreatedAt",
						Type: &GoType{
							Name: "time.Time",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &GenerateGoProcess{}
			got, err := p.GetGoStruct(tt.args.messageProto, tt.args.messageIndex, tt.args.protoFileInfo, tt.args.protoTypeInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateGoProcess.GetGoStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Hogefunc differs: (-got +want)\n%s", diff)
			}
		})
	}
}
