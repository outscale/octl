# Usage overview

## General syntax

```sh
octl <command> <subcommand> [flags]
````

## Main commands

- `iaas` - Core IaaS API
- `kube` - OKS API
- `profile` - Profile management
- `update` - Update to the latest version
- `completion` - Generate completion shell script

## High-level commands vs API calls
High-level commands provide opinionated defaults (like tables) and shortcuts:

- `octl iaas <entity> list`
- `octl iaas <entity> describe <id> [<id>]...`
- `octl iaas <entity> create`
- `octl iaas <entity> update <id> [<id>]... <flags>`
- `octl iaas <entity> delete <id> [<id>]...`

Direct API calls are also available:

- `octl iaas api <OperationName> <flags>`

## Output
Output-related flags are documented here:
[outputs.md](./outputs.md)

## Flag syntax
The flag syntax is:

- list of values are comma-separated: `--Filters.VmStateNames running,stopped`
- boolean flags can be set to false by setting: `--TmpEnabled=false`
- lists of embedded objects (e.g. `Nics` or `BlockDeviceMappings` in `CreateVms`) can be configured using indexes:
    - `--BlockDeviceMappings.0.Bsu.VolumeType`

- time flag values can be set:
    - using the RFC3339 format (e.g. `2026-02-10T14:52:30Z`)
    - as a duration offset with a `+` or `-` prefix (e.g. `+10m`, `-1h`)
    - as a day/month/year offset with a `+` or `-` prefix (e.g. `+1mo`, `-1y`)

## Next
- IaaS examples:    [iaas.md](./iaas.md)
- Output formats:   [outputs.md](./outputs.md)
- jq and filters:   [jq-and-filters.md](./jq-and-filters.md)
- Templating:       [templating.md](./templating.md)
- Chaining:         [chaining.md](./chaining.md)
- Columns and expr: [columns-and-expr.md](./columns-and-expr.md)