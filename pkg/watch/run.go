package watch

import (
	"context"
	"errors"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/outscale/octl/pkg/messages"
	"github.com/outscale/octl/pkg/output/format"
	"github.com/outscale/octl/pkg/spinner"
)

type RefreshMsg struct{}

// Run runs the watcher.
func Run[Error error](ctx context.Context, fmter format.Interface, fn func(ctx context.Context) error, interval time.Duration) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	cancel = spinner.Run(ctx, fmt.Sprintf("Refreshing every %s... Press ctrl+c to exit", interval))
	defer cancel()
	model, ok := fmter.(*Format)
	if !ok {
		messages.ExitErr(errors.New("invalid output for --watch"))
	}
	prog := tea.NewProgram(model)

	// Start the API calling loop
	exitErr := make(chan error, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := fn(ctx)
				if err != nil {
					exitErr <- err
					prog.Quit()
					return
				}
				// prog.Send(RefreshMsg{})
				// The cursor sometimes reappears... This also forces a refresh.
				prog.Send(tea.HideCursor())
				time.Sleep(interval)
			}
		}
	}()

	// Start Bubbletea
	if _, err := prog.Run(); err != nil {
		return fmt.Errorf("watch: %w", err)
	}

	// Forward the error from the API loop
	select {
	case err := <-exitErr:
		return err
	default:
		return nil
	}
}
