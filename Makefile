prodmigrateup:
	migrate -database "$(DB_URL)" -path db/migrations up

prodmigratedown:
	migrate -database "$(DB_URL)" -path db/migrations down

stagingmigrateup:
	migrate -database "$(DB_STAGING_URL)" -path db/migrations up

stagingmigratedown:
	migrate -database "$(DB_STAGING_URL)" -path db/migrations down
