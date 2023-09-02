updb:
	docker-compose -f docker-compose.yml up --build
downdb:
	docker-compose -f docker-compose.yml down
killdb:
	docker-compose -f docker-compose.yml down --volumes
runapp:
	go run cmd/main.go
testdb:
	go test ./internal/repository -cover
