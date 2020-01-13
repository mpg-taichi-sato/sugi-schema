package generator

import (
	"context"
	"strings"
)

type GoFile struct {
	PackageName    string
	PackageComment string
	Structs        []*GoStruct
}

type GoStruct struct {
	Name    string
	Fields  []*GoStructField
	Comment string
}

type GoStructField struct {
	Name    string
	Type    *GoType
	Comment string
}

type GoType struct {
	Name string
}

type GoCodeGenerator struct {
}

func (g *GoCodeGenerator) Generate(ctx context.Context, f *GoFile) string {
	var sb strings.Builder
	// canvas := &Canvas{}

	if f.PackageComment != "" {
		sb.WriteString(g.GoComment(f.PackageComment))
	}

	// package
	packageName := f.PackageName
	sb.WriteString("package ")
	sb.WriteString(packageName)

	// import
	importTime := false
	for _, st := range f.Structs {
		for _, field := range st.Fields {
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
	for i := 0; i < len(f.Structs); i++ {
		st := f.Structs[i]
		sb.WriteString("\n\n")
		if st.Comment != "" {
			sb.WriteString("//")
			sb.WriteString(st.Name)
			sb.WriteString(" ")
			sb.WriteString(st.Comment)
		}
		sb.WriteString("type ")
		sb.WriteString(st.Name)
		sb.WriteString(" struct {\n")
		for j := 0; j < len(st.Fields); j++ {
			field := st.Fields[j]
			sb.WriteString("\t")
			sb.WriteString(field.Name)
			sb.WriteString(" ")
			sb.WriteString(field.Type.Name)

			if field.Comment != "" {
				sb.WriteString(" //")
				sb.WriteString(field.Comment)
			} else {
				sb.WriteString("\n")
			}
		}
		sb.WriteString("}")
	}

	return sb.String()
}

func (g *GoCodeGenerator) GoComment(src string) string {
	var sb strings.Builder
	lines := strings.Split(src, "\n")
	for i := 0; i < len(lines); i++ {
		if i == len(lines)-1 && lines[i] == "" {
			break
		}
		sb.WriteString("//")
		sb.WriteString(lines[i])
		sb.WriteString("\n")
	}
	return sb.String()
}
