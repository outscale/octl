package flags

import (
	"fmt"

	"github.com/spf13/pflag"
)

const aliasNoForward = "octl.outscale.com/no-forward"

func MarkAsNoForward(fs *pflag.FlagSet, name string) error {
	return fs.SetAnnotation(name, aliasNoForward, []string{""})
}

func IsNoForward(f *pflag.Flag) bool {
	_, found := f.Annotations[aliasNoForward]
	return found
}

const defaultValue = "octl.outscale.com/default"

func SetDefault(fs *pflag.FlagSet, name, value string) error {
	f := fs.Lookup(name)
	if f == nil {
		return fmt.Errorf("flag %q not found", name)
	}
	SetDefaultValue(f, value)
	return nil
}

func SetDefaultValue(f *pflag.Flag, value string) {
	f.DefValue = value
	_ = f.Value.Set(value)
	if f.Annotations == nil {
		f.Annotations = map[string][]string{}
	}
	f.Annotations[defaultValue] = []string{value}
}

func HasDefault(f *pflag.Flag) bool {
	val, found := f.Annotations[defaultValue]
	return found && len(val) > 0
}
