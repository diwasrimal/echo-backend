DB = echochat-db

.PHONY: db test server clean

all: server db

server:
	go build .

db:
	createdb $(DB)
	psql -d $(DB) -f ./db/sql/create_tables.sql

dbrecreate:
	psql -d $(DB) -f ./db/sql/drop_tables.sql
	psql -d $(DB) -f ./db/sql/create_tables.sql

dbrefill: dbrecreate
	psql -d $(DB) -f ./db/sql/empty_tables.sql
	psql -d $(DB) -f ./db/sql/fill_tables.sql

test:
	go test -v ./...

clean:
	go clean
	dropdb $(DB)
