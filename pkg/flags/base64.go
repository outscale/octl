package flags

import (
	"encoding/base64"
	"errors"
	"os"

	"github.com/outscale/octl/pkg/debug"
)

const Base64File = "base64File"

var ErrInvalidBase64 = errors.New("value is neither a file nor base64 data")

// Base64FileValue adapts File.File for use as a flag.
type Base64FileValue struct {
	content []byte
}

func NewBase64FileValue() *Base64FileValue {
	return &Base64FileValue{}
}

// Set sets the value based on a file content.
func (v *Base64FileValue) Set(s string) error {
	if _, err := os.Stat(s); err == nil {
		v.content, err = os.ReadFile(s) //nolint:gosec
		return err
	}
	content, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return ErrInvalidBase64
	}
	v.content = content
	return nil
}

// Type name for File.File flags.
func (v *Base64FileValue) Type() string {
	return "base64File"
}

func (v *Base64FileValue) String() string {
	if len(v.content) == 0 {
		return ""
	}
	_, err := base64.StdEncoding.DecodeString(string(v.content))
	if err == nil {
		debug.Println("value is already base 64 encoded")
		return string(v.content)
	}
	return base64.StdEncoding.EncodeToString(v.content)
}

func (v *Base64FileValue) Value() (string, bool) {
	return v.String(), true
}
