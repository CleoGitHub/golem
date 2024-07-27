package router

import "context"

type TokenValidator interface {
	ValidateAtLeastOneRole(ctx context.Context, token string, roles []string) (bool, error)
}
