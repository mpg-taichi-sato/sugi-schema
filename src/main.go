// Flagを読み取ってどのprocessを用いるのかを決定する
package main

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"protoc-gen-genta/generator"
	"strings"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

func parseReq(r io.Reader) (*plugin.CodeGeneratorRequest, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var req plugin.CodeGeneratorRequest
	if err = proto.Unmarshal(buf, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func parseRequestParameter(parameter string) *generator.Options {
	options := &generator.Options{}
	for _, p := range strings.Split(parameter, ",") {
		spec := strings.SplitN(p, "=", 2)
		if len(spec) == 1 {
			switch spec[0] {
			case "go":
				options.GenGo = true
			case "go_tag_json":
				options.GoTagJSON = true
			case "json":
				options.GenJSON = true
			case "apidoc":
				options.GenAPIDoc = true
			case "csfields":
				options.GenCSFields = true
			}
			continue
		}
	}
	return options
}

func emitResp(resp *plugin.CodeGeneratorResponse) error {
	buf, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(buf)
	return err
}

type GenerateProcess interface {
	Run(ctx context.Context, req *plugin.CodeGeneratorRequest) ([]*plugin.CodeGeneratorResponse_File, error)
}

func run() error {
	// request取得
	req, err := parseReq(os.Stdin)
	if err != nil {
		return err
	}

	if req.Parameter == nil {
		return errors.New("no parameter")
	}

	options := parseRequestParameter(req.GetParameter())

	// response生成
	var resp plugin.CodeGeneratorResponse

	var processes []GenerateProcess
	if options.GenGo {
		processes = append(processes, &generator.GenerateGoProcess{
			TagJSON: options.GoTagJSON,
		})
	}

	if options.GenJSON {
		processes = append(processes, &generator.GenerateJSONProcess{})
	}

	if options.GenAPIDoc {
		processes = append(processes, &generator.GenerateAPIDocProcess{})
	}

	if options.GenCSFields {
		processes = append(processes, &generator.GenerateCSFieldsProcess{})
	}

	if len(processes) == 0 {
		return errors.New("no processes")
	}
	ctx := context.Background()
	for _, process := range processes {
		responseFiles, err := process.Run(ctx, req)
		if err != nil {
			return err
		}
		if len(responseFiles) == 0 {
			continue
		}

		resp.File = append(resp.File, responseFiles...)
	}

	return emitResp(&resp)
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
