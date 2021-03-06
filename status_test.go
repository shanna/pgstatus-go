package pgstatus_test

import (
	"testing"

	"github.com/shanna/pgstatus-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type pgerror string

func (p pgerror) Error() string {
	return string(p)
}

func (p pgerror) SQLState() string {
	return string(p)
}

func TestConvert(t *testing.T) {
	testCases := []struct {
		in   error
		want *status.Status
	}{
		{in: nil, want: status.New(codes.OK, "")},
		{in: status.New(codes.DeadlineExceeded, "deadline exceeded").Err(), want: status.New(codes.DeadlineExceeded, "deadline exceeded")},
		{in: pgerror("00000"), want: status.New(codes.OK, "successful_completion: 00000")},
		{in: pgerror("0100C"), want: status.New(codes.OK, "dynamic_result_sets_returned: 0100C")},
		{in: pgerror("03000"), want: status.New(codes.Unavailable, "sql_statement_not_yet_complete: 03000")},
		{in: pgerror("23505"), want: status.New(codes.AlreadyExists, "unique_violation: 23505")},
	}
	for _, tc := range testCases {
		got := pgstatus.Convert(tc.in)
		if got.Code() != tc.want.Code() {
			t.Errorf("FromContextError(%v) = %v; want %v", tc.in, got.Code(), tc.want.Code())
		}

		if got.Message() != tc.want.Message() {
			t.Errorf("FromContextError(%v) = %s; want %s", tc.in, got.Message(), tc.want.Message())
		}
	}
}

func TestCode(t *testing.T) {
	testCases := []struct {
		in   error
		want codes.Code
	}{
		{in: nil, want: codes.OK},
		{in: status.New(codes.DeadlineExceeded, "deadline exceeded").Err(), want: codes.DeadlineExceeded},
		{in: pgerror("00000"), want: codes.OK},
		{in: pgerror("0100C"), want: codes.OK},
		{in: pgerror("03000"), want: codes.Unavailable},
		{in: pgerror("23505"), want: codes.AlreadyExists},
	}
	for _, tc := range testCases {
		got := pgstatus.Code(tc.in)
		if got != tc.want {
			t.Errorf("FromContextError(%v) = %v; want %v", tc.in, got, tc.want)
		}
	}
}
