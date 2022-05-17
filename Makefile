postgres:
	docker run --name messagedb -p 5588:5432  -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -d postgres:14.2-alpine

createdb:
	docker exec -it messagedb createdb --username=root --owner=root message_db

dropdb:
	docker exec -it messagedb dropdb message_db

upmigrate:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5588/message_db?sslmode=disable" -verbose up

dnmigrate:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5588/message_db?sslmode=disable" -verbose down

.PHONY: postgres createdb dropdb upmigrate