# go-boilerplate
boilerplate for go API project

### DB Design

https://dbdiagram.io/d/5fec2f3a80d742080a3494e7

### Migration
 
migrate boilerplatetest db down

go run cmd/dbmigrate/main.go \
-migrate=down \
-dbname=boilerplatetest \
-dbhost=localhost

you can now repeat steps from above to connect to pg container
and ensure that users table is missing from boilerplatetest DB.

now bring it back up

go run cmd/dbmigrate/main.go \
-migrate=up \
-dbname=boilerplatetest \
-dbhost=localhost

### Connect to PSQL Docker

docker exec -it go-boilerplate_pg_1 psql -U postgres -d boilerplatetest

### Connect to Redis

docker exec -it go-boilerplate_redis_1 redis-cli

### Mock Store using gomock package

mockgen -package mockstore -destination internal/store/mock/store.go github.com/amryamanah/go-boilerplate/internal/store/sqlc Store
