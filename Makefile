binary_name=app
config=config/default.yaml

.PHONY:
	build \
	run \
	up \
	down \

build:
	go build -o $(binary_name) cmd/api-server/*.go

run: build
	./$(binary_name) -config=$(config)

up:
	sudo docker compose up -d

down:
	sudo docker compose down
