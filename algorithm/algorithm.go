package algorithm

import (
	"context"

	algorithmio_pb "github.com/d-sparks/gravy/algorithm/proto"

	"google.golang.org/grpc"
)

// An A is an interface that all algorithm grpc clients should implement.
type A interface {
	Execute(
		ctx context.Context,
		in *algorithmio_pb.Input,
		opts ...grpc.CallOption,
	) (
		*algorithmio_pb.Output,
		error,
	)
}
