build:
	@go build -o httx cmd/tcplistener/main.go
	@./httx

run:
	@go run cmd/tcplistener/main.go

clean:
	@rm httx
