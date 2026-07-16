package cmd

import (
	"context"
	"strings"

	"github.com/outscale/goutils/sdk/metadata"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/octl/pkg/messages"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func deviceToVolumeID(cmd *cobra.Command, cl *osc.Client) error {
	if err := doDeviceToVolumeID(cmd, cl); err != nil {
		return err
	}
	return doDevicesToVolumeIDs(cmd, cl)
}

func doDeviceToVolumeID(cmd *cobra.Command, cl *osc.Client) error {
	f := cmd.Flag("VolumeId")
	if f == nil {
		return nil
	}
	dev := f.Value.String()
	if !strings.HasPrefix(dev, "/dev/") {
		return nil
	}
	vols, err := fetchVolumeIDs(cmd.Context(), cl, dev)
	switch {
	case err != nil:
		return err
	case len(vols) > 0:
		_ = f.Value.Set(vols[0])
	}
	return nil
}

func doDevicesToVolumeIDs(cmd *cobra.Command, cl *osc.Client) error {
	f := cmd.Flag("Filters.VolumeIds")
	if f == nil {
		return nil
	}
	fs, ok := f.Value.(pflag.SliceValue)
	if !ok {
		return nil
	}
	devs := lo.Filter(fs.GetSlice(), func(dev string, _ int) bool {
		return strings.HasPrefix(dev, "/dev/")
	})
	if len(devs) == 0 {
		return nil
	}
	vols, err := fetchVolumeIDs(cmd.Context(), cl, devs...)
	switch {
	case err != nil:
		return err
	case len(vols) > 0:
		keep := lo.Filter(fs.GetSlice(), func(dev string, _ int) bool {
			return !strings.HasPrefix(dev, "/dev/")
		})
		keep = append(keep, vols...)
		_ = fs.Replace(keep)
	}
	return nil
}

func fetchVolumeIDs(ctx context.Context, cl *osc.Client, devs ...string) ([]string, error) {
	vm, err := metadata.GetInstanceID(ctx)
	if err != nil {
		return nil, err
	}
	messages.Info("Fetching device(s) %s from instance %s", strings.Join(devs, ", "), vm)
	resp, err := cl.ReadVolumes(ctx, osc.ReadVolumesRequest{
		Filters: &osc.FiltersVolume{
			LinkVolumeDeviceNames: &devs,
			LinkVolumeVmIds:       &[]string{vm},
		},
	})
	switch {
	case err != nil:
		return nil, err
	case len(ptr.From(resp.Volumes)) > 0:
		vols := lo.Map(*resp.Volumes, func(v osc.Volume, _ int) string {
			return v.VolumeId
		})
		messages.Info("Using volume ID(s) %s", strings.Join(vols, ", "))
		return vols, nil
	default:
		return nil, nil
	}
}
