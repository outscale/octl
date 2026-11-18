package format_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/outscale/octl/pkg/output/format"
	"github.com/outscale/osc-sdk-go/v3/pkg/oks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYAML_OneOfFields(t *testing.T) {
	item := oks.ErrorItem{
		Code: "422",
		Type: "error",
	}
	_ = item.Details.FromErrorItemDetails1(oks.ErrorItemDetails1{
		{Msg: "bad"},
	})
	v := oks.ErrorResponse{
		Errors: []oks.ErrorItem{item},
	}
	fmt := format.NewYAML("")
	buf := bytes.Buffer{}
	err := fmt.Format(t.Context(), &buf, v)
	require.NoError(t, err)
	assert.Equal(t, `Errors:
- Code: '422'
  Details:
  - loc: null
    msg: bad
    type: ''
  Type: error
ResponseContext: {}`, strings.TrimSpace(buf.String()))
}
