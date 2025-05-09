package middleware_grpc

import (
	"context"

	apperrors "github.com/SamEkb/messenger-app/auth-service/internal/app/errors"
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
			case apperrors.Is(err, apperrors.ErrNotFound):
				code = codes.NotFound
			case apperrors.Is(err, apperrors.ErrAlreadyExists):
				code = codes.AlreadyExists
			case apperrors.Is(err, apperrors.ErrUnauthorized):
				code = codes.Unauthenticated
			case apperrors.Is(err, apperrors.ErrInvalidInput),
				apperrors.Is(err, apperrors.ErrValidation):
				code = codes.InvalidArgument
			case apperrors.Is(err, apperrors.ErrDatabaseConnection),
				apperrors.Is(err, apperrors.ErrDatabaseQuery):
				code = codes.Unavailable
			case apperrors.Is(err, apperrors.ErrTimeout):
				code = codes.DeadlineExceeded
			default:
				code = codes.Internal
			}

			return nil, status.Error(code, err.Error())
		}

		return resp, nil
	}
}
