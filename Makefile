run:
	go run cmd/main.go

swag:
	swag init -g internal/http/handler/handler.go