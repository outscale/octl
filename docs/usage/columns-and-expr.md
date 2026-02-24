# Changing table columns and expr language

## Changing table columns

Columns can be replaced:

```sh
octl iaas vm list --columns "ID:VmId|DNS:PrivateDnsName"
┌────────────┬───────────────────────────────────────────┐
│     ID     │                    DNS                    │
├────────────┼───────────────────────────────────────────┤
│ i-foo      │ ip-10-1-112-23.eu-west-2.compute.internal │
│ i-bar      │ ip-10-9-35-211.eu-west-2.compute.internal │
│ i-baz      │ ip-10-0-4-143.eu-west-2.compute.internal  │
└────────────┴───────────────────────────────────────────┘
```

Columns can be added to the standard columns:

```sh
octl iaas vm list --columns +DNS:PrivateDnsName
```

## expr language

Column content is defined with the [expr language](https://expr-lang.org/docs/language-definition).

To display a tag value:

`octl iaas vm list --columns "+tag:find(Tags, #?.Key == \"Name\")?.Value"`