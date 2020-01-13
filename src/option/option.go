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
