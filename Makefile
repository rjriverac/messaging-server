DB_URL=postgresql://root:secret@localhost:5588/message_db?sslmode=disable

postgres:
	docker run --name messagedb --network message-network -p 5588:5432  -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -d postgres:14.2-alpine

createdb:
	docker exec -it messagedb createdb --username=root --owner=root message_db

dropdb:
	docker exec -it messagedb dropdb message_db

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

upmigrate:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

downnmigrate:
	migrate -path db/migration -database "$(DB_URL)" -verbose down
	
migratedown1:	
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc-dev generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/rjriverac/messaging-server/db/sqlc Store

db_docs:
	dbdocs build docs/db.dbml

db_schema:
	dbml2sql --postgres -o docs/schema.sql docs/db.dbml

.PHONY: postgres db_docs db_schema createdb dropdb upmigrate downnmigrate sqlc test server mock