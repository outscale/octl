package cmd_test

import (
	"testing"

	"github.com/outscale/goutils/sdk/metadata"
	"github.com/stretchr/testify/assert"
)

func TestMetadata(t *testing.T) {
	t.Run("metadata all returns all metadata", func(t *testing.T) {
		var meta metadata.Metadata
		runJSON(t, []string{"metadata", "all", "-o", "json"}, nil, &meta)
		assert.NotEmpty(t, meta.Hostname)
		assert.NotEmpty(t, meta.InstanceID)
		assert.NotEmpty(t, meta.InstanceType)
		assert.NotEmpty(t, meta.OMIID)
		assert.NotEmpty(t, meta.MAC)
		assert.NotEmpty(t, meta.Placement.Subregion)
		assert.NotEmpty(t, meta.Placement.Cluster)
		assert.NotEmpty(t, meta.Placement.Server)
		assert.NotEmpty(t, meta.DeviceMapping)
	})
	t.Run("metadata hostname returns a result", func(t *testing.T) {
		res := run(t, []string{"metadata", "hostname"}, nil)
		assert.NotEmpty(t, res)
	})
	t.Run("metadata availability_zone returns a result", func(t *testing.T) {
		res := run(t, []string{"metadata", "availability_zone"}, nil)
		assert.NotEmpty(t, res)
	})
	t.Run("metadata region returns a result", func(t *testing.T) {
		res := run(t, []string{"metadata", "region"}, nil)
		assert.NotEmpty(t, res)
	})
	t.Run("metadata instance_id returns a result", func(t *testing.T) {
		res := run(t, []string{"metadata", "instance_id"}, nil)
		assert.NotEmpty(t, res)
	})
	t.Run("metadata instance_type returns a result", func(t *testing.T) {
		res := run(t, []string{"metadata", "instance_type"}, nil)
		assert.NotEmpty(t, res)
	})
	t.Run("metadata mac returns a result", func(t *testing.T) {
		res := run(t, []string{"metadata", "mac"}, nil)
		assert.NotEmpty(t, res)
	})
	t.Run("metadata placement_cluster returns a result", func(t *testing.T) {
		res := run(t, []string{"metadata", "placement_cluster"}, nil)
		assert.NotEmpty(t, res)
	})
	t.Run("metadata placement_server returns a result", func(t *testing.T) {
		res := run(t, []string{"metadata", "placement_server"}, nil)
		assert.NotEmpty(t, res)
	})
	t.Run("metadata device_mapping returns a result", func(t *testing.T) {
		var res map[string]string
		runJSON(t, []string{"metadata", "device_mapping", "-o", "json"}, nil, &res)
		assert.NotEmpty(t, res)
		assert.Equal(t, "/dev/sda1", res["root"])
	})
}
