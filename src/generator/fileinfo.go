package generator

import (
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// ProtoFileInfo pathを意識せずに使えるようにする
type ProtoFileInfo struct {
	sourceCodeInfoLocations []*descriptor.SourceCodeInfo_Location
}

func CreateProtoFileInfo(fileDescriptorProto *descriptor.FileDescriptorProto) *ProtoFileInfo {
	sourceCodeInfo := fileDescriptorProto.GetSourceCodeInfo()
	sourceCodeInfoLocations := sourceCodeInfo.GetLocation() // ロケーションの配列
	return &ProtoFileInfo{
		sourceCodeInfoLocations: sourceCodeInfoLocations,
	}
}

func (info *ProtoFileInfo) GetSyntaxInfo() *descriptor.SourceCodeInfo_Location {
	if info.sourceCodeInfoLocations == nil {
		return nil
	}
	for i := range info.sourceCodeInfoLocations {
		loc := info.sourceCodeInfoLocations[i]
		path := loc.Path
		if len(path) != 1 {
			continue
		}

		// field syntax: 12
		if path[0] != 12 {
			continue
		}

		return loc
	}
	return nil
}

func (info *ProtoFileInfo) GetPackageInfo() *descriptor.SourceCodeInfo_Location {
	if info.sourceCodeInfoLocations == nil {
		return nil
	}
	for i := range info.sourceCodeInfoLocations {
		loc := info.sourceCodeInfoLocations[i]
		path := loc.Path
		if len(path) != 1 {
			continue
		}

		// field package: 2
		if path[0] != 2 {
			continue
		}

		return loc
	}
	return nil
}

func (info *ProtoFileInfo) GetMethodInfo(serviceIndex, methodIndex int32) *descriptor.SourceCodeInfo_Location {
	if info.sourceCodeInfoLocations == nil {
		return nil
	}
	for i := range info.sourceCodeInfoLocations {
		loc := info.sourceCodeInfoLocations[i]
		path := loc.Path
		if len(path) != 4 {
			continue
		}

		// field service: 6 method 2
		if path[0] != 6 || path[1] != serviceIndex || path[2] != 2 || path[3] != methodIndex {
			continue
		}

		return loc
	}
	return nil
}

func (info *ProtoFileInfo) GetServiceInfo(serviceIndex int32) *descriptor.SourceCodeInfo_Location {
	if info.sourceCodeInfoLocations == nil {
		return nil
	}
	for i := range info.sourceCodeInfoLocations {
		loc := info.sourceCodeInfoLocations[i]
		path := loc.Path
		if len(path) != 2 {
			continue
		}

		// field service: 6
		if path[0] != 6 || path[1] != serviceIndex {
			continue
		}

		return loc
	}
	return nil
}

func (info *ProtoFileInfo) GetMessageInfo(messageIndex int32) *descriptor.SourceCodeInfo_Location {
	if info.sourceCodeInfoLocations == nil {
		return nil
	}
	for i := range info.sourceCodeInfoLocations {
		loc := info.sourceCodeInfoLocations[i]
		path := loc.Path
		if len(path) != 2 {
			continue
		}

		// field message: 4
		if path[0] != 4 || path[1] != messageIndex {
			continue
		}

		return loc
	}
	return nil
}

func (info *ProtoFileInfo) GetMessageFieldInfo(messageIndex, fieldIndex int32) *descriptor.SourceCodeInfo_Location {
	if info.sourceCodeInfoLocations == nil {
		return nil
	}
	for i := range info.sourceCodeInfoLocations {
		loc := info.sourceCodeInfoLocations[i]
		path := loc.Path
		if len(path) != 4 {
			continue
		}

		// field message: 4 field: 2
		if path[0] != 4 || path[1] != messageIndex || path[2] != 2 || path[3] != fieldIndex {
			continue
		}

		return loc
	}
	return nil
}
