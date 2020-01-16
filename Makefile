generate:
	rm -rf build/gen
	mkdir -p build/gen
	protoc --proto_path=pb --plugin=src/protoc-gen-genta --genta_out=go,json:build/gen pb/*.proto
	protoc --proto_path=pb --plugin=src/protoc-gen-genta --genta_out=go,go_tag_json,json,apidoc:build/gen pb/api/*.proto
	protoc --proto_path=pb --plugin=src/protoc-gen-genta --genta_out=go,go_tag_json,json,csfields:build/gen pb/masterdata/*.proto
	go fmt ./build/gen/...

generate-option:
	protoc --proto_path=pb --go_out=src pb/option/*.proto