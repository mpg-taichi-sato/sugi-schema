package option

import (
	"github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func GetAPIOption(options *descriptor.MethodOptions) *Http {
	if options == nil {
		return nil
	}
	ext, err := proto.GetExtension(options, E_Http)
	if err == proto.ErrMissingExtension {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	apiOption := ext.(*Http)
	return apiOption
}

func GetGoTagOption(options *descriptor.FieldOptions) (string, bool) {
	if options == nil {
		return "", false
	}
	ext, err := proto.GetExtension(options, E_GoTag)
	if err == proto.ErrMissingExtension {
		return "", false
	}
	if err != nil {
		panic(err)
	}
	tag, ok := ext.(*string)
	if !ok {
		return "", false
	}
	return *tag, true
}
