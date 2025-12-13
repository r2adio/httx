build:
	@go build
	@./httx

run:
	@go run main.go

clean:
	@rm httx
