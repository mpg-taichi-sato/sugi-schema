generate:
	mkdir -p build/gen
	protoc --proto_path=pb --plugin=src/protoc-gen-genta --genta_out=go,json:build/gen pb/*.proto 
	protoc --proto_path=pb --plugin=src/protoc-gen-genta --genta_out=json:build/gen pb/api/*.proto
	protoc --proto_path=pb --plugin=src/protoc-gen-genta --genta_out=apidoc:build/gen pb/*.proto pb/api/*.proto
	go fmt ./build/gen/...

generate-option:
	protoc --proto_path=pb --go_out=src pb/option/*.proto