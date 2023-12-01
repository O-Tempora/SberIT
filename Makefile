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
	sudo docker run -d -p 5555:5432 --name=db -e POSTGRES_DB=sber -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres postgres
	timeout 15s bash -c 'until sudo docker exec db pg_isready ; do sleep 1 ; done'
	./$(binary_name) -config=$(config)

up:
	sudo docker compose up -d

down:
	sudo docker compose down
