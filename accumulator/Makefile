
run:
	go run main.go

mocks:
	rm -rfd mocks
	go generate ./...

test: mocks
	go test -timeout 500ms ./...

integration_test:
	docker-compose down
	docker-compose up -d --build
	sleep 5
	go test --tags integration --timeout=30s -v ./...
	docker-compose down