package flags

import (
	"os"
)

const File = "file"

// FileValue sets a flag with either a file content.
type FileValue struct {
	content []byte
}

func NewFileValue() *FileValue {
	return &FileValue{}
}

// Set sets the value based on a file content.
func (v *FileValue) Set(s string) error {
	if _, err := os.Stat(s); err == nil {
		v.content, err = os.ReadFile(s) //nolint:gosec
		return err
	}
	// fallback for retro-compatibility
	v.content = []byte(s)
	return nil
}

// Type name for File.File flags.
func (v *FileValue) Type() string {
	return "file"
}

func (v *FileValue) String() string {
	return string(v.content)
}

func (v *FileValue) Value() (string, bool) {
	return v.String(), true
}
