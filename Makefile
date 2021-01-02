sqlc:
	sqlc generate
test:
	go test -v -cover ./...
mock:
	mockgen -package mockstore -destination internal/store/mock/store.go github.com/amryamanah/go-boilerplate/internal/store/sqlc Store

.PHONY: sqlc test mock