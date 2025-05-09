package middleware_grpc

import (
	"context"
	"fmt"

	"github.com/SamEkb/messenger-app/auth-service/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorsUnaryServerInterceptor преобразует ошибки приложения в gRPC ошибки
func ErrorsUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			// Если это уже gRPC статус, возвращаем как есть
			if _, ok := status.FromError(err); ok {
				return resp, err
			}

			// Преобразуем ошибки приложения в gRPC ошибки
			var code codes.Code
			switch {
			case errors.Is(err, errors.ErrNotFound):
				code = codes.NotFound
			case errors.Is(err, errors.ErrAlreadyExists):
				code = codes.AlreadyExists
			case errors.Is(err, errors.ErrUnauthorized),
				errors.Is(err, errors.ErrTokenExpired),
				errors.Is(err, errors.ErrInvalidToken):
				code = codes.Unauthenticated
			case errors.Is(err, errors.ErrForbidden):
				code = codes.PermissionDenied
			case errors.Is(err, errors.ErrInvalidInput),
				errors.Is(err, errors.ErrValidation):
				code = codes.InvalidArgument
			case errors.Is(err, errors.ErrDatabaseConnection),
				errors.Is(err, errors.ErrDatabaseQuery),
				errors.Is(err, errors.ErrServiceUnavailable):
				code = codes.Unavailable
			default:
				code = codes.Internal
			}

			// Получаем детали ошибки для более информативного сообщения
			details := errors.GetErrorDetails(err)
			statusErr := status.Error(code, err.Error())

			// Если есть детали, добавляем их в статус
			if details != nil && len(details) > 0 {
				// Для простоты, просто добавляем строковое представление деталей в статусе
				detailsStr := ""
				for k, v := range details {
					detailsStr += fmt.Sprintf("%s: %v; ", k, v)
				}
				if detailsStr != "" {
					return nil, status.Errorf(code, "%s [Details: %s]", err.Error(), detailsStr)
				}
			}

			return nil, statusErr
		}

		return resp, nil
	}
}
