## Echo Backend

This is the backend code repository for [echo](https://github.com/diwasrimal/echo) chat application.

### Build

#### Requirements
* Go
* PostgreSQL

```console
git clone https://github.com/diwasrimal/echo-backend
cd echo-backend
createdb echochat-db                                # -+
psql -d echochat-db -f ./db/sql/create_tables.sql   #  |- or just `make`
go build .                                          # -+
```

### Run

Make your own `.env` from `.env.example` with correct values and run.

```console
./echo-backend
```
