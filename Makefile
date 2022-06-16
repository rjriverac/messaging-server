postgres:
	docker run --name messagedb --network message-network -p 5588:5432  -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -d postgres:14.2-alpine

createdb:
	docker exec -it messagedb createdb --username=root --owner=root message_db

dropdb:
	docker exec -it messagedb dropdb message_db

upmigrate:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5588/message_db?sslmode=disable" -verbose up

downnmigrate:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5588/message_db?sslmode=disable" -verbose down

sqlc:
	sqlc-dev generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/rjriverac/messaging-server/db/sqlc Store

.PHONY: postgres createdb dropdb upmigrate downnmigrate sqlc test server mock