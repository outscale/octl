package spinner

import (
	"context"
	"os"

	"github.com/charmbracelet/huh/spinner"
	"github.com/outscale/octl/pkg/style"
)

func Run(ctx context.Context, text string) (cancel func()) {
	spinWait := make(chan struct{})
	spinCtx, spinCancel := context.WithCancel(ctx)
	spin := spinner.New().
		Title(text).
		Context(spinCtx).
		Output(os.Stderr).
		Style(style.Yellow).
		TitleStyle(style.Faint)
	go func() {
		_ = spin.Run()
		close(spinWait)
	}()

	return func() {
		spinCancel()
		<-spinWait
	}
}
