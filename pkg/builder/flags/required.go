package flags

import (
	"strings"
)

const required = "[REQUIRED]"

func RequiredHelp(help string) string {
	if strings.HasPrefix(help, required) {
		return help
	}
	return required + " " + help
}
