migrate:
	goose -dir internal/database/migrations sqlite3 gpsitty.db up
run:
	go run ./cmd/api/main.go