package handlers

type StateGettable interface {
	GetState(userId int64) string
}
