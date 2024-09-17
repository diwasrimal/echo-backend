db := "echochat-db"

default: server db

server:
	go build .

db:
	createdb {{db}}
	psql -d {{db}} -f ./db/sql/create_tables.sql

dbrecreate:
	psql -d {{db}} -f ./db/sql/drop_tables.sql
	psql -d {{db}} -f ./db/sql/create_tables.sql

dbrefill: dbrecreate
	psql -d {{db}} -f ./db/sql/empty_tables.sql
	psql -d {{db}} -f ./db/sql/fill_tables.sql

test:
	go test -v ./...

clean:
	go clean
	dropdb {{db}}
