package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Lysander66/ace/pkg/crpc/_examples/pb"
)

var _ pb.AddServer = (*AddSvc)(nil)

type AddSvc struct {
	pb.UnimplementedAddServer
}

const (
	intMax = 1<<31 - 1
	intMin = -(intMax + 1)
	maxLen = 10
)

var (
	// ErrTwoZeroes is an arbitrary business rule for the Add method.
	ErrTwoZeroes = errors.New("can't sum two zeroes")

	// ErrIntOverflow protects the Add method. We've decided that this error
	// indicates a misbehaving service and should count against e.g. circuit
	// breakers. So, we return it directly in endpoints, to illustrate the
	// difference. In a real service, this probably wouldn't be the case.
	ErrIntOverflow = errors.New("integer overflow")

	// ErrMaxSizeExceeded protects the Concat method.
	ErrMaxSizeExceeded = errors.New("result exceeds maximum size")
)

func (s AddSvc) Sum(ctx context.Context, req *pb.SumRequest) (*pb.SumReply, error) {
	slog.Info("Sum", "a", req.A, "b", req.B)
	a, b := req.A, req.B
	if a == 0 && b == 0 {
		return nil, ErrTwoZeroes
	}
	if (b > 0 && a > (intMax-b)) || (b < 0 && a < (intMin-b)) {
		return nil, ErrIntOverflow
	}
	return &pb.SumReply{V: a + b}, nil
}

func (s AddSvc) Concat(ctx context.Context, req *pb.ConcatRequest) (*pb.ConcatReply, error) {
	slog.Info("Concat", "a", req.A, "b", req.B)
	a, b := req.A, req.B
	if len(a)+len(b) > maxLen {
		return nil, ErrMaxSizeExceeded
	}
	return &pb.ConcatReply{V: a + b}, nil
}
