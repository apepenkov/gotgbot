package events

type Event interface {
	ImplementsEvent()
	RollbackActions() error
	ExecuteActions() error
}
