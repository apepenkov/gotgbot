package handlers

import "context"

type StateGettable interface {
	GetState(userId int64, ctx context.Context) (string, error)
}
