## octl kube api GetNodepoolTemplate

Retrieves the default node pool template, including predefined configurations for node scaling, storage, and upgrade strategies.

### Synopsis

Retrieves the default node pool template, including predefined configurations for node scaling, storage, and upgrade strategies. Returns a JSON response containing the request context and template details, such as node count, node type, availability zones, volume specifications, upgrade strategy, and auto-healing settings.

```
octl kube api GetNodepoolTemplate [flags]
```

### Options

```
  -h, --help   help for GetNodepoolTemplate
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

* [octl kube api](octl_kube_api.md)	 - kube api calls

