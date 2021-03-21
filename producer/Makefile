
run:
	go run main.go

mocks:
	rm -rfd mocks
	go generate ./...

test: mocks
	go test --timeout=5s ./...

integration_test:
	docker-compose up -d --build
	go test --tags integration --timeout=5s ./...
	docker-compose down