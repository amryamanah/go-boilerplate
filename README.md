# go-boilerplate
boilerplate for go API project


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
