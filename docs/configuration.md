# Configuration

`octl` expects either environment variables or a profile file.

## Environment variables

The tool will try to read:
- `OSC_ACCESS_KEY`
- `OSC_SECRET_KEY`
- `OSC_REGION`

## Profile file

If no environment variables are defined, the tool expects to find a profile in a profile file.

- Default config path: `~/.osc/config.json`
  - Can be set with `--config` or `OSC_CONFIG_FILE`
- Default profile name: `default`
  - Can be set with `--profile` or `OSC_PROFILE`

If `--config` or `--profile` is set, `octl` will load the profile file without checking environment variables.

### Example

```json
{
  "default": {
    "access_key": "MyAccessKey",
    "secret_key": "MySecretKey",
    "region": "eu-west-2"
  }
}
```
## Profile management

`octl profile` allows you to manage your profile file:

- octl profile list
- octl profile add
- octl profile delete
- octl profile use

> Note: Environment variables take precedence. A profile marked as the default may not be used if relevant environment variables are set.