[![](https://godoc.org/github.com/shanna/pgstatus-go?status.svg)](https://godoc.org/github.com/shanna/pgstatus-go)

# pgstatus-go

Postgres error to gRPC status.

Supports any `error` implementing a `SQLState() string` method that:
1. Returns a [postgres error code](https://www.postgresql.org/docs/current/static/errcodes-appendix.html).
2. Isn't already a gRPC `*status.Status` pointer.
3. Doesn't implement a `GRPCStatus() *grpc.Status`. 

## Drivers

### github.com/jackc/pgx

Supported as of `github.com/jackc/pgx/v4`. YMMV if you are using older version that don't use 
`github.com/jackc/pgconn` as I don't know the history of the `SQLState() string` interface
in the driver.

### github.com/lib/pq

Unsupported.

The pq package provides no usable interface. Casting to a *pq.Error and providing a
`SQLState() string` is currently an exercise left to the reader.

### github.com/go-pg/pg

Unsupported.

## License

MIT

