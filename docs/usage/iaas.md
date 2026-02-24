# IaaS usage

## High level commands

High level commands are available for the IaaS provider.

### List

`octl iaas <entity> list` lists all entities using the `table` format:

```sh
octl iaas volume list
┌──────────────┬──────────┬───────────┬───────────────┬──────┬──────┐
│      ID      │   Type   │   State   │ SubregionName │ Size │ Iops │
├──────────────┼──────────┼───────────┼───────────────┼──────┼──────┤
│ vol-foo      │ io1      │ in-use    │ eu-west-2a    │ 300  │ 5000 │
│ vol-bar      │ standard │ in-use    │ eu-west-2a    │ 20   │ 0    │
│ vol-baz      │ gp2      │ available │ eu-west-2a    │ 4    │ 100  │
└──────────────┴──────────┴───────────┴───────────────┴──────┴──────┘
```

### Describe

`octl iaas <entity> describe <id> [<id>]...` displays one or multiple entities using the `yaml` format:

```sh
octl iaas volume describe vol-foo
CreationDate: '2024-12-17T11:07:58.757Z'
Iops: 5000
LinkedVolumes:
- DeviceName: /dev/sda1
  State: attached
  VmId: i-foo
  VolumeId: vol-foo
Size: 300
SnapshotId: snap-foo
State: in-use
SubregionName: eu-west-2a
Tags: []
VolumeId: vol-foo
VolumeType: io1
```

### Create

```sh
octl iaas <entity> create creates an entity:

octl iaas vol create --subregion-name eu-west-2a --size 4
CreationDate: '2026-02-19T15:37:47.015Z'
LinkedVolumes: []
Size: 4
State: creating
SubregionName: eu-west-2a
Tags: []
VolumeId: vol-foo
VolumeType: standard
```

### Update

`octl iaas <entity> update <id> [<id>]... <flags>` updates one or multiple entities with the same parameters:

```sh
octl iaas vol update vol-foo vol-bar --size 6
```

### Delete

`octl iaas <entity> delete <id> [<id>]...` deletes one or multiple entities:

```sh
octl iaas vol delete vol-foo vol-bar
```

## API access

The API can be directly called, with a `raw` output:

```sh
octl iaas api ReadVms --Filters.VmStateNames running
```