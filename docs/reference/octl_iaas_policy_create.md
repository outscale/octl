## octl iaas policy create

alias for api CreatePolicy

### Synopsis

> *alias for api CreatePolicy*

Creates a managed policy to apply to a user.

This action creates a policy version and sets v1 as the default one.

```
octl iaas policy create [flags]
```

### Options

```
      --description string    A description for the policy.
      --document fileOrJson   [REQUIRED] Either a file storing the policy document or the policy document (in JSON format).
  -h, --help                  help for create
      --name string           [REQUIRED] The name of the policy.
      --path string           The path of the policy.
```

### Options inherited from parent commands

```
      --color-by string             color lines by this item (JQ syntax)
  -c, --columns string              columns to display - [+]<title>:<jq query for content>||<title>:<jq query for content>
      --config string               Path of profile file (by default, ~/.osc/config.json)
      --dry-run                     Display the request payload that would be sent to the API without sending it
      --filter strings              comma separated list of filters for results - name:value,name:value, alias for jq filter 'select(.name | tostring | test("value"))'
      --jq string                   jq filter
      --no-upgrade                  do not check for new versions
  -O, --out-file string             redirect output to file
  -o, --output string               output format (raw, json, yaml, table, csv, none, base64, text)
      --payload string              JSON content for query body
      --profile string              Profile to use in profile file (by default, "default")
      --single                      convert single entry lists to a single object
      --template string             JSON template file for query body
  -v, --verbose                     Verbose output
      --waitfor string              jq expression to wait for - octl will query every waitfor-interval until the expression returns 1/true or a non empty result
      --waitfor-interval duration   interval between two waitfor iterations (default 5s)
      --waitfor-timeout duration    maximum duration of a wait (default 10m0s)
  -y, --yes                         answer yes to all prompts
```

### SEE ALSO

* [octl iaas policy](octl_iaas_policy.md)	 - policy commands

