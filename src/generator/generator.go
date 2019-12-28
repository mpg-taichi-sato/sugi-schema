package generator

import (
	"strings"
)

type File struct {
	PackageName  string
	StructValues []*StructValue
}

type StructValue struct {
	Name   string
	Fields []*StructField
}

type StructField struct {
	Name string
	Type *Type
}

type Type struct {
	Name string
}

func GenerateGoCode(f *File) string {
	var sb strings.Builder
	// canvas := &Canvas{}
	// package
	packageName := f.PackageName
	sb.WriteString("package ")
	sb.WriteString(packageName)

	// import
	importTime := false
	for _, structValue := range f.StructValues {
		for _, field := range structValue.Fields {
			if field.Type.Name == "time.Time" {
				importTime = true
				break
			}
		}
		if importTime {
			break
		}
	}
	if importTime {
		sb.WriteString("\n\n")
		sb.WriteString("import (\n")
		sb.WriteString("\t\"time\"\n")
		sb.WriteString(")\n")
	}

	// message
	for i := 0; i < len(f.StructValues); i++ {
		structValue := f.StructValues[i]
		sb.WriteString("\n\n")

		sb.WriteString("type ")
		sb.WriteString(structValue.Name)
		sb.WriteString(" struct {\n")
		for j := 0; j < len(structValue.Fields); j++ {
			field := structValue.Fields[j]
			sb.WriteString("\t")
			sb.WriteString(field.Name)
			sb.WriteString(" ")
			sb.WriteString(field.Type.Name)
			sb.WriteString("\n")
		}
		sb.WriteString("}")
	}

	return sb.String()
}
