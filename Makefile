migrate:
	goose -dir internal/database/migrations sqlite3 gpsitty.db up