server: sqlc
	go run main.go server

migrate:
	go run main.go migrate

mock:
	mockgen -source pkg/identity/merchant/domain/repository/repository.go -destination internal/mock/gen_repo_merchant.go -package mock

sqlc:
	sqlc compile && sqlc generate

mc:
	migrate create -ext sql -seq -digits 3 -dir ./internal/db/migrations $(n)

migrate.rollback:
	migrate -path ./db/migrations -database "mysql://root:abcd1234@tcp(localhost:3306)/cs-api?charset=utf8mb4&multiStatements=true&parseTime=true" -verbose down $(n)

test:
	go test ./internal/module/...

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server main.go