// Package interceptors provides various middleware functionality for GRPC.
package interceptors

import (
	"context"
	"dk-go-gophkeeper/internal/config"
	"dk-go-gophkeeper/internal/server/cipher"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthHandler defines attributes and methods of an AuthHandler instance.
type AuthHandler struct {
	sec cipher.Cipher
	cfg *config.Config
}

// NewAuthHandler initializes AuthHandler instance.
func NewAuthHandler(sec cipher.Cipher, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		sec: sec,
		cfg: cfg,
	}
}

// AuthFunc checks request context for metadata and validates metadata-derived authorization token.
func (a *AuthHandler) AuthFunc(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "No authorization data was found")
	}
	values := md.Get(a.cfg.AuthBearerName)
	if len(values) == 0 {
		return status.Error(codes.Unauthenticated, "Empty authorization data was found")
	}
	_, err := a.sec.Decode(values[0])
	if err != nil {
		return status.Error(codes.PermissionDenied, err.Error())
	}
	return nil
}

// UnaryServerInterceptor returns a new unary server interceptors that performs per-request authentication.
func (a *AuthHandler) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		switch info.FullMethod {
		case "/proto.Gophkeeper/Login":
			return handler(ctx, req)
		case "/proto.Gophkeeper/Register":
			return handler(ctx, req)
		default:
			err := a.AuthFunc(ctx)
			if err != nil {
				return nil, err
			}
			return handler(ctx, req)
		}
	}
}
