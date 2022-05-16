postgres:
	docker run --name messagedb -p 5588:5432  -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -d postgres:14.2-alpine

createdb:
	docker exec -it messagedb createdb --username=root --owner=root message_db

dropdb:
	docker exec -it messagedb dropdb message_db
	
.PHONY: postgres createdb dropdb