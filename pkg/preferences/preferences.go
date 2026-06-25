package preferences

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/outscale/octl/pkg/debug"
	"github.com/outscale/octl/pkg/messages"
)

type Kube struct {
	DefaultProject string `yaml:"default_project,omitempty"`
}

type All struct {
	Kube Kube `yaml:"kube,omitzero"`
}

var Preferences All

func init() {
	root, err := os.UserConfigDir()
	if err != nil {
		debug.Println("Unable to compute user preferences dir", err)
		return
	}
	path := filepath.Join(root, "octl", "preferences.yaml")
	fd, err := os.Open(path) //nolint:gosec
	if errors.Is(err, fs.ErrNotExist) {
		debug.Println("no user config file found", path)
		return
	}
	if err != nil {
		messages.Err("Unable to open preferences path: %w")
		return
	}
	debug.Println("loading user config from", path)
	err = yaml.NewDecoder(fd).Decode(&Preferences)
	if err != nil {
		messages.Err("Unable to load preferences: %w")
	}
}

func (p *All) Save() (err error) {
	root, err := os.UserConfigDir()
	if err != nil {
		debug.Println("Unable to compute user preferences dir", err)
		return err
	}
	dir := filepath.Join(root, "octl")
	err = os.MkdirAll(dir, 0o700)
	if err != nil {
		return err
	}
	path := filepath.Join(dir, "preferences.yaml")
	fd, err := os.Create(path) //nolint:gosec
	if err != nil {
		return err
	}
	defer func() {
		cerr := fd.Close()
		if err == nil {
			err = cerr
		}
	}()
	debug.Println("saving user preferences to", path)
	return yaml.NewEncoder(fd).Encode(p)
}
