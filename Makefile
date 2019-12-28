generate:
	mkdir -p build/gen
	protoc --proto_path=pb --plugin=src/protoc-gen-model --model_out=build/gen pb/*.proto
	go fmt ./build/gen/...