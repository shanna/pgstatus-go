// Package pgstatus converts postgres errors to gRPC statuses.
package pgstatus

import (
	"database/sql"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SQLState error code interface.
type SQLState interface {
	SQLState() string
}

// FromError returns a gRPC status.Status representing an err if it was produced from status package, has a method
// `GRPCStatus() *status.Status` or has a method SQLState interface returning a postgres error code.
// Otherwise, ok is false and a Status is returned with codes.Unknown and the original error message.
func FromError(err error) (*status.Status, bool) {
	if err == nil {
		return nil, true
	}

	if err == sql.ErrNoRows {
		return status.New(codes.NotFound, err.Error()), true
	}

	if se, ok := status.FromError(err); ok {
		return se, true
	}

	if condition, ok := ConditionName(err); ok {
		return status.New(Code(err), fmt.Sprintf("%s: %s", condition, err.Error())), true
	}

	return status.New(codes.Unknown, err.Error()), false
}

// Convert is a convenience function which removes the need to handle the boolean return value from FromError.
func Convert(err error) *status.Status {
	s, _ := FromError(err)
	return s
}

// Code converts a database error into a gRPC codes.Code.
//
// Postgres: https://www.postgresql.org/docs/current/static/errcodes-appendix.html
// GRPC:     https://godoc.org/google.golang.org/grpc/codes
func Code(err error) codes.Code {
	if err == nil {
		return codes.OK
	}

	if err == sql.ErrNoRows {
		return codes.NotFound
	}

	if se, ok := status.FromError(err); ok {
		return se.Code()
	}

	se, ok := err.(SQLState)
	if !ok {
		return codes.Unknown
	}

	// Pull request for postgres to gRPC code changes accompanied by reasoned opinions will probably be accepted.
	code := se.SQLState()
	switch {
	// cool
	case code == "00000":
		return codes.OK

	// warnings only.
	case strings.HasPrefix(code, "01"):
		return codes.OK

	// no data
	case strings.HasPrefix(code, "02"):
		return codes.NotFound

	// not complete
	case strings.HasPrefix(code, "03"):
		return codes.Unavailable

	// connection error
	case strings.HasPrefix(code, "08"):
		return codes.Unavailable

	// triggered action exception
	case strings.HasPrefix(code, "09"):
		return codes.Internal

	// invalid grantor
	case strings.HasPrefix(code, "0L"):
		return codes.PermissionDenied

	// invalid role specification
	case strings.HasPrefix(code, "0P"):
		return codes.PermissionDenied

	// foreign key violation
	case code == "23503":
		return codes.FailedPrecondition

	// uniqueness violation
	case code == "23505":
		return codes.AlreadyExists

	// invalid transaction state
	case strings.HasPrefix(code, "25"):
		return codes.Aborted

	// invalid auth specification
	case strings.HasPrefix(code, "28"):
		return codes.PermissionDenied

	// invalid transaction termination
	case strings.HasPrefix(code, "2D"):
		return codes.Internal

	// external routine exception
	case strings.HasPrefix(code, "38"):
		return codes.Internal

	// external routine invocation
	case strings.HasPrefix(code, "39"):
		return codes.Internal

	// savepoint exception
	case strings.HasPrefix(code, "3B"):
		return codes.Aborted

	// transaction rollback
	case strings.HasPrefix(code, "40"):
		return codes.Aborted

	// syntax errors or access rule violations
	case strings.HasPrefix(code, "42"):
		return codes.Internal

	// insufficient resources
	case strings.HasPrefix(code, "53"):
		return codes.ResourceExhausted

	// too complex
	case strings.HasPrefix(code, "54"):
		return codes.Internal

	// obj not in prereq state
	case strings.HasPrefix(code, "55"):
		return codes.Internal

	// operator intervention
	case strings.HasPrefix(code, "57"):
		return codes.Internal

	// system error
	case strings.HasPrefix(code, "58"):
		return codes.Internal

	// conf file error
	case strings.HasPrefix(code, "F0"):
		return codes.Internal

	// TODO(shane): Foreign data wrapper errors?

	// default code for "raise"
	case code == "P0001":
		return codes.Aborted

	// PL/pgSQL error
	case strings.HasPrefix(code, "P0"):
		return codes.Internal

	// internal error
	case strings.HasPrefix(code, "XX"):
		return codes.Internal

	// undefined function
	case code == "42883":
		return codes.Internal

	// undefined table
	case code == "42P01":
		return codes.Internal

	// insufficient privileges
	case code == "42501":
		return codes.PermissionDenied

	default:
		return codes.Internal
	}
}
