binary_name=app
config=config/default.yaml

.PHONY:
	build \
	run \

build:
	go build -o $(binary_name) cmd/server/*.go

run: build
	./$(binary_name) -config=$(config)
