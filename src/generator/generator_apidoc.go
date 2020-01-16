package generator

import (
	"context"
	"fmt"
	"strings"
)

type APIDocInfo struct {
	HeadComment string
	Services    []*APIDocService
	SubStructs  []*APIDocStructModel
}

type APIDocService struct {
	Method   string
	Path     string
	Comment  string
	Request  *APIDocStructModel
	Response *APIDocStructModel
}

type APIDocStructModel struct {
	Name    string
	Fields  []*APIDocStructField
	Comment string
}

type APIDocStructField struct {
	Name        string
	Description string
	DataType    string
}

type APIDocGenerator struct {
}

func (g *APIDocGenerator) Generate(ctx context.Context, docInfo *APIDocInfo) string {
	var sb strings.Builder
	sb.WriteString(docInfo.HeadComment)

	for i := 0; i < len(docInfo.Services); i++ {
		service := docInfo.Services[i]
		sb.WriteString(fmt.Sprintf("# %s %s\n", service.Method, service.Path))
		sb.WriteString(service.Comment)
		if service.Request != nil {
			sb.WriteString("##### Parameters  \n")
			sb.WriteString("|Parameter|Description|Data Type|\n")
			sb.WriteString("|:--|:--|:--|\n")
			for j := 0; j < len(service.Request.Fields); j++ {
				parameter := service.Request.Fields[j]
				sb.WriteString(fmt.Sprintf("|%s|%s|%s|\n", parameter.Name, parameter.Description, parameter.DataType))
			}
		}
		if service.Response != nil {
			sb.WriteString("##### Response  \n")
			sb.WriteString("|Parameter|Description|Data Type|\n")
			sb.WriteString("|:--|:--|:--|\n")
			for j := 0; j < len(service.Response.Fields); j++ {
				parameter := service.Response.Fields[j]
				sb.WriteString(fmt.Sprintf("|%s|%s|%s|\n", parameter.Name, parameter.Description, parameter.DataType))
			}
		}
	}

	if len(docInfo.SubStructs) > 0 {
		sb.WriteString("\n\n")
		sb.WriteString("# SubStructs\n")
	}

	for i := 0; i < len(docInfo.SubStructs); i++ {
		s := docInfo.SubStructs[i]
		sb.WriteString(fmt.Sprintf("#### %s\n", s.Name))
		sb.WriteString(s.Comment)
		sb.WriteString("|Parameter|Description|Data Type|\n")
		sb.WriteString("|:--|:--|:--|\n")
		for j := 0; j < len(s.Fields); j++ {
			parameter := s.Fields[j]
			sb.WriteString(fmt.Sprintf("|%s|%s|%s|\n", parameter.Name, parameter.Description, parameter.DataType))
		}
	}

	return sb.String()
}
