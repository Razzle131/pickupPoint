.SILENT:

run:
	go run cmd/main.go

unit-test:
	go test ./... -cover -short

integr-test-full-output:
	go test -v -run Integration

integr-test:
	go test -run Integration

image:=foo
docker-build:
	docker build -t $(image) .

docker-run:
	docker run --rm -p 8080:8080 --env-file .env $(image)

lint:
	golangci-lint run