# Troubleshooting

## `octl update` fails

`octl update` requires write access to the binary location.

- If the binary is in a protected directory, run with appropriate permissions.
- Otherwise, download the latest binary from:
  - https://github.com/outscale/octl/releases/latest

## Completion does not work

Verify you installed the completion script for your shell and that it is sourced:

- Bash: ensure the file is placed in `/etc/bash_completion.d/` and reload your shell
- Zsh: ensure `_octl` is in a directory listed in `$fpath` (commonly `/usr/local/share/zsh/site-functions`)
- Fish: ensure the file is in `/usr/share/fish/completions/`

## Table output falls back to YAML

When using `--jq` with `-o table`, `octl` may not be able to format the jq result into a table.
In that case it will fall back to YAML output.

Example:

```sh
octl iaas volume list --jq '.VolumeId' -o table
```

## Profile not used

Environment variables take precedence over profiles. If `OSC_ACCESS_KEY`, `OSC_SECRET_KEY`, and `OSC_REGION`
are set, the profile file may be ignored unless you explicitly force a config/profile selection.