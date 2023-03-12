test:
	go test ./... -v -race -cover -coverprofile=coverage.txt && go tool cover -func=coverage.txt

lint:
	golangci-lint run --deadline=5m -v

build:
	go build -o ./bin/server ./cmd/server/
	go build -o ./bin/client ./cmd/client/

docker_build:
	docker build --target=client -t pow-client -f docker/Dockerfile .
	docker build --target=server -t pow-server -f docker/Dockerfile .

run_client:
	docker run --env SERVER_HOST=$(SERVER_HOST) pow-client

run_server:
	docker run -p 8081:8081 pow-server
