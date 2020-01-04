generate:
	mkdir -p build/gen
	protoc --proto_path=pb --plugin=src/protoc-gen-model --model_out=build/gen pb/*.proto pb/api/*.proto
	go fmt ./build/gen/...

generate-option:
	protoc --proto_path=pb --go_out=src pb/option/*.proto