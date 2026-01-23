.PHONY: run generate migrate-up migrate-down

run:
	go build cmd/server/main.go && go run cmd/server/main.go

gen:
	cd internal && go run github.com/99designs/gqlgen generate

dev:
	go tool air

install-tools:
	go install github.com/99designs/gqlgen@latest
	go install github.com/air-verse/air@latest