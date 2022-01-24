# TO CREATE A MIGRATION, FROM PROJECT ROOT FOLDER:
# migrate create -ext sql -dir db/migrations -seq <MIGRATIONNAME>
# migrate -path db/migrations -database <dburl> version
# migrate -path db/migrations -database <dburl> force <n>
# migrate -path db/migrations -database <dburl> up
# migrate -path db/migrations -database <dburl> down
# in case of dirty output: DELETE FROM schema_migrations WHERE 1=1; 

prodmigrateup:
	migrate -database "$(DB_URL)" -path db/migrations up

prodmigratedown:
	migrate -database "$(DB_URL)" -path db/migrations down

stagingmigrateup:
	migrate -database "$(DB_STAGING_URL)" -path db/migrations up

stagingmigratedown:
	migrate -database "$(DB_STAGING_URL)" -path db/migrations down
