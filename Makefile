obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu
receiver:
	@go build -o bin/receiver ./data_receiver
	@./bin/receiver
up:
	@docker compose up -d

.PHONY: obu
