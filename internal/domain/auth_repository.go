package domain

import (
	"context"
)

type AuthRepository interface {
	GetIdAndHash(ctx context.Context, email string) (string, string, error)
}
