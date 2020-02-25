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
		{in: pgerror("00000"), want: status.New(codes.OK, "00000")},
		{in: pgerror("0100C"), want: status.New(codes.OK, "0100C")},
		{in: pgerror("03000"), want: status.New(codes.Unavailable, "03000")},
		{in: pgerror("23505"), want: status.New(codes.AlreadyExists, "23505")},
	}
	for _, tc := range testCases {
		got := pgstatus.Convert(tc.in)
		if got.Code() != tc.want.Code() || got.Message() != tc.want.Message() {
			t.Errorf("FromContextError(%v) = %v; want %v", tc.in, got, tc.want)
		}
	}
}
