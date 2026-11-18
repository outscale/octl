package structs_test

import (
	"testing"

	"github.com/outscale/octl/pkg/structs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToUnstructured(t *testing.T) {
	v := struct {
		Foo string `json:"foo"`
	}{
		Foo: "bar",
	}
	vv, err := structs.ToUnstructured(v)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"foo": "bar"}, vv)
	vvv, err := structs.ToUnstructured(vv)
	require.NoError(t, err)
	assert.Equal(t, vv, vvv)
}
