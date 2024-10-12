package cb_data

import "strings"

type CallbackData string

const Sep = "\x01"

func BuildCallbackData(data CallbackData, args ...string) string {
	return strings.Join(append([]string{string(data)}, args...), Sep)
}

func GetCallbackData(data string) (CallbackData, []string) {
	parts := strings.Split(data, Sep)
	if len(parts) == 0 {
		return "", nil
	}

	return CallbackData(parts[0]), parts[1:]
}

var BCB = BuildCallbackData
