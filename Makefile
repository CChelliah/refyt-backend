migrateup:
	migrate -path libs/database/migration -database "postgresql://cavinashchelliah:cavinashchelliah@refyt.crxkhygkhnbk.ap-southeast-2.rds.amazonaws.com:5432/postgres?sslmode=disable" -verbose up
migratedown:
	migrate -path libs/database/migration -database "postgresql://cavinashchelliah:password@localhost:5432/postgres?sslmode=disable" -verbose down