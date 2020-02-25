# pgstatus-go

Convert postgres error codes to HTTP and gRPC status codes.

## Drivers

Supports any `error` implementing a `SQLState() string` method that:
1. Returns a [postgres error code](https://www.postgresql.org/docs/current/static/errcodes-appendix.html).
2. Isn't already a gRPC `*status.Status` pointer.
3. Doesn't implement a `GRPCStatus() *grpc.Status`. 

### github.com/jackc/pgx

Supported. 

### github.com/lib/pq

No support.

The pq package provides no usable interface.
Casting to a *pq.Error and providing a `SQLState() string` is currently left to the reader.

## License

MIT

