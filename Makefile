updb:
	docker-compose -f docker-compose.yml up --build
downdb:
	docker-compose -f docker-compose.yml down
killdb:
	docker-compose -f docker-compose.yml down --volumes
app:
	go run cmd/main.go
testdb:
	go test ./internal/repository -cover
detailed_testdb:
	go test ./internal/repository -v -cover
