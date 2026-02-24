# jq and filters

## Using jq filters

Based on raw payload (`-o raw`):

```sh
octl iaas api ReadVms --jq ".Vms[].VmId"
```

Based on content (`-o json`/`-o yaml`/`-o table`):

```sh
octl iaas api ReadVms --jq ".VmId" -o json
```

or

```sh
octl iaas vm list --jq ".VmId" -o json
```

`--jq` will try to output to tables if possible:

```sh
octl iaas volume list --jq 'select(.State | test("in-use"))'
┌──────────────┬──────┬──────────┬───────────┬──────┬──────┬────────────┬───────────┐
│      ID      │ Name │   Type   │   State   │ Size │ Iops │     VM     │  Device   │
├──────────────┼──────┼──────────┼───────────┼──────┼──────┼────────────┼───────────┤
│ vol-foo      │      │ io1      │ in-use    │ 300  │ 5000 │ i-foo      │ /dev/sda1 │
│ vol-bar      │      │ standard │ in-use    │ 20   │      │ i-bar      │ /dev/sda1 │
└──────────────┴──────┴──────────┴───────────┴──────┴──────┴────────────┴───────────┘
```

but:
```sh
octl iaas volume list --jq '.VolumeId' -o table
Unable to format as a table, switching to YAML...
- vol-foo
- vol-bar
````

## Using filters (`--filter`)

With `--filter`, a list of content filters can be set.

To display the list of images for Kubernetes v1.31 for example:

`octl iaas image list --filter ImageName:kubernetes,ImageName:v1.31`

This is the equivalent of running the two following jq filters:

`select(.ImageName | test("kubernetes"))`

`select(.ImageName | test("v1.31"))`