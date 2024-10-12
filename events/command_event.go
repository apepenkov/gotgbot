package events

import "strings"

type CommandEvent struct {
	NewMessageEvent
	Command      string
	MentionedBot string
	Arguments    []string
}

func (e *CommandEvent) ImplementsEvent() {}

func NewCommandEvent(event NewMessageEvent) *CommandEvent {
	parts := strings.Split(event.Message.Text, " ")
	var command, mentionedBot string
	var arguments []string

	if strings.Contains(parts[0], "@") {
		mentionParts := strings.Split(parts[0], "@")
		command = mentionParts[0]
		mentionedBot = mentionParts[1]
	} else {
		command = parts[0]
	}

	command = strings.TrimPrefix(command, "/")

	if len(parts) > 1 {
		arguments = parts[1:]
	} else {
		arguments = []string{}
	}

	return &CommandEvent{
		NewMessageEvent: event,
		Command:         command,
		MentionedBot:    mentionedBot,
		Arguments:       arguments,
	}
}
