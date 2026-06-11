# Working with self-signed certificates

> The following configuration is **unsafe** and must only be used in development and never in a production setup.

You can configure octl to skip TLS certificate validation.

When using environment-based authentication, you can use the `OSC_TLS_SKIP_VERIFY` environment variable:
```shell
export OSC_ACCESS_KEY=<ak>
export OSC_SECRET_KEY=<sk>

OSC_TLS_SKIP_VERIFY=true octl iaas vm ls
```

When using a profile file, you can add the `tls_skip_verify` key to the profile:
```json
{
  "<profile name>": {
    "access_key": "<ak>",
    "secret_key": "<sk>",
    "region": "<region>",
    "tls_skip_verify": true
  }
}
```
