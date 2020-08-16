

grpc-build-builder:
	docker build -f grpc_builder.Dockerfile -t mg-grpc-builder .

grpc-generate:
	docker run \
		-v $(shell pwd):/workdir \
		--workdir /workdir \
		mg-grpc-builder \
		\
		protoc \
		-I ./ $(shell find . -name "*.proto") \
		--go_out=plugins=grpc:. \
		--go_opt=paths=source_relative

grpc-lint:
	# lint protocol buffer files
	docker run --rm --volume "$(shell pwd):/workspace" --workdir /workspace bufbuild/buf check lint
	# load latest remote master
	git fetch
	# check breaking changes
	docker run --rm --volume "$(shell pwd):/workspace" --workdir /workspace bufbuild/buf check breaking --against-input '.git#branch=origin/master'
