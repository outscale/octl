## octl iaas securitygrouprule delete

Deletes one or more inbound or outbound rules from a security group.

### Synopsis

Deletes one or more inbound or outbound rules from a security group.

For the rule to be deleted, the values specified in the deletion request must exactly match the value of the existing rule.

In case of TCP and UDP protocols, you have to indicate the destination port or range of ports. In case of ICMP protocol, you have to specify the ICMP type and code numbers.

Rules (IP permissions) consist of the protocol, IP range or source security group.

To remove outbound access to a destination security group, we recommend to use a set of IP permissions. We also recommend to specify the protocol in a set of IP permissions.


Alternatively, you can use the `Rules` parameter to delete several rules at the same time.

> alias for DeleteSecurityGroupRule

```
octl iaas securitygrouprule delete [flags]
```

### Options

```
      --flow string                         [REQUIRED] The direction of the flow: Inbound or Outbound. (default "Inbound")
      --group-id string                     [REQUIRED] The ID of the security group you want to delete a rule from.
  -h, --help                                help for delete
      --ports strings                       A list of either protocol (all ports from a protocol, e.g. icmp), protocol/port (a single port/protocol, e.g. tcp/80) or protocol/from-to (a range, e.g. tcp/8080-8082)
      --remote-account string               The OUTSCALE account ID that owns the source or destination security group.
      --remote-ranges strings               One or more IP ranges for the security group rules, in CIDR notation (for example, ["10.0.0.0/24" , "10.0.1.0/24"]).
      --remote-security-group string        The ID of a source or destination security group that you want to link to the security group of the rule.
      --remote-security-group-name string   The name of a source or destination security group that you want to link to the security group of the rule.
      --remote-service strings              One or more service IDs to allow traffic from a Net to access the corresponding OUTSCALE services.
```

### Options inherited from parent commands

```
  -c, --columns string             columns to display - [+]<title>:<jq query for content>||<title>:<jq query for content>
      --config string              Path of profile file (by default, ~/.osc/config.json)
      --dry-run                    Display the request payload that would be sent to the API without sending it
      --elapsed                    add elapsed time column when using --watch (default true)
      --filter strings             comma separated list of filters for results - name:value,name:value, alias for jq filter 'select(.name | tostring | test("value"))'
      --interval duration          interval between two watch/waitfor iterations (default 5s)
      --jq string                  jq filter
      --max-pages int              maximum number of pages a command can fetch (default 20)
      --no-upgrade                 do not check for new versions
  -O, --out-file string            redirect output to file
  -o, --output string              output format (json, yaml, raw, rawyaml, table, csv, none, text)
      --payload string             JSON content for query body
      --profile string             Profile to use in profile file (by default, "default")
  -s, --silent                     Hides all information messages
      --single                     convert single entry lists to a single object
      --style string               style to use for syntax-highlighting (doom-one, github, monokai, nord, paraiso, solarized) (default "github")
      --template string            JSON template file for query body
  -v, --verbose                    Verbose output
      --waitfor string             repeatedly call the API until the specified jq expression returns 1/true or a non empty result
      --waitfor-timeout duration   maximum duration of a wait (default 10m0s)
      --watch                      repeatedly call the API and display changes
  -y, --yes                        answer yes to all prompts
```

### SEE ALSO

* [octl iaas securitygrouprule](octl_iaas_securitygrouprule.md)	 - Manage SecurityGroupRule resources

