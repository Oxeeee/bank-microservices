run:
	go run ./cmd/billing --config=./configs/local.yaml

migrate:
	go run ./cmd/migrator --dsn=postgres://postgres:fingo123@localhost:5434?sslmode=disable --migrations-path=file://migrations
