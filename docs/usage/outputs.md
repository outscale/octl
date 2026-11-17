# Output formats and output-related options

## Output-related options

| Option | Default | Allowed values | Description |
| ------ | ------- | -------------- | ----------- |
| `--jq` | | | jq-like filter |
| `--filter` | | | content filter |
| `--template` | | | JSON template for query body |
| `--output` | | `json`, `yaml`, `table`, `csv`, `base64`, `text`, `none` | output format |
| `--raw` | | | raw output (as returned by the API without parsing & pagination) |
| `--columns` | | | columns to display in a table |
| `--single` | | | convert single entry lists to a single object |

## Output formats

- `json` displays the content in JSON format, without response context,
- `yaml` displays the content in YAML format, without response context,
- `table` displays the content in CSV format, based on columns defined by the `--columns` flag,
- `csv` displays the content in a text table, based on columns defined by the `--columns` flag,
- `base64` decodes base64-encoded strings or lists of strings,
- `text` displays the content as plain text,
- `none` disables output.

> Please note that when using `--raw`, octl returns the raw payload instead of only outputting the content
(e.g. a list of VMs when listing VMs instead of an object with a `Vms` attribute storing the list).

## Pagination

Unless `--raw` is set, octl will repeatedly fetch the API to get all pages of content (with a maximum of `--max-pages` pages, 20 by default).

