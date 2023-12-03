binary_name=app
db_container_name=db
config=config/default.yaml

.PHONY:
	build \
	run \
	up \
	down \
	killdb \
	test \

build:
	go build -o $(binary_name) cmd/api-server/*.go

run: build
	sudo docker run -d -p 5555:5432 --name=$(db_container_name) -e POSTGRES_DB=sber -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres postgres
	timeout 15s bash -c 'until sudo docker exec $(db_container_name) pg_isready ; do sleep 1 ; done'
	./$(binary_name) -config=$(config)

killdb:
	sudo docker stop $(db_container_name) && sudo docker rm $(db_container_name)

up:
	sudo docker compose up -d

down:
	sudo docker compose down

test:
	sudo docker run -d -p 6969:5432 --name=testdb -e POSTGRES_DB=sber -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres postgres
	timeout 15s bash -c 'until sudo docker exec testdb pg_isready ; do sleep 1 ; done'
	go test -count=1 ./...
	sudo docker stop testdb && sudo docker rm testdb
