package spinner

import (
	"context"
	"io"
	"os"

	"github.com/charmbracelet/huh/spinner"
	"github.com/outscale/octl/pkg/style"
)

var Out io.Writer = os.Stderr

func Run(ctx context.Context, text string, foreground bool) (cancel func()) {
	spinWait := make(chan struct{})
	spinCtx, spinCancel := context.WithCancel(ctx)
	spin := spinner.New().
		Title(text).
		Context(spinCtx).
		Output(Out).
		Style(style.Yellow).
		TitleStyle(style.Faint)
	if foreground {
		_ = spin.Run()
		spinCancel()
		return
	}
	go func() {
		_ = spin.Run()
		close(spinWait)
	}()

	return func() {
		spinCancel()
		<-spinWait
	}
}
