package middleware_grpc

import (
	"context"

	"github.com/SamEkb/messenger-app/friends-service/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorsUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			if _, ok := status.FromError(err); ok {
				return resp, err
			}

			var code codes.Code
			switch {
			case errors.Is(err, errors.ErrNotFound):
				code = codes.NotFound
			case errors.Is(err, errors.ErrAlreadyExists):
				code = codes.AlreadyExists
			case errors.Is(err, errors.ErrUnauthorized):
				code = codes.Unauthenticated
			case errors.Is(err, errors.ErrForbidden):
				code = codes.PermissionDenied
			case errors.Is(err, errors.ErrInvalidInput),
				errors.Is(err, errors.ErrValidation):
				code = codes.InvalidArgument
			case errors.Is(err, errors.ErrDatabaseConnection),
				errors.Is(err, errors.ErrDatabaseQuery):
				code = codes.Unavailable
			case errors.Is(err, errors.ErrTimeout):
				code = codes.DeadlineExceeded
			default:
				code = codes.Internal
			}

			return nil, status.Error(code, err.Error())
		}

		return resp, nil
	}
}
