migrateup:
	migrate -path libs/database/migration -database "postgresql://cavinashchelliah:password@localhost:5432/postgres?sslmode=disable" -verbose up

migratedown:
	migrate -path libs/database/migration -database "postgresql://cavinashchelliah:password@localhost:5432/postgres?sslmode=disable" -verbose down