package events

type CommandEvent struct {
	NewMessageEvent
	Command      string
	MentionedBot string
	Arguments    []string
}

func (e *CommandEvent) ImplementsEvent() {}
