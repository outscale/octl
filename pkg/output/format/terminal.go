package format

import (
	"io"
	"os"

	"github.com/mattn/go-isatty"
)

func IsTerminal(w io.Writer) bool {
	if fd, ok := w.(*os.File); ok {
		return isatty.IsTerminal(fd.Fd())
	}
	return false
}
