package domain

import "context"

type PLDRepository interface {
	IsValidUser(ctx context.Context, user *User) (bool, error)
}
