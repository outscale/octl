# Sending raw JSON and templating

## Sending raw JSON

```sh
echo '{"SubnetId":"subnet-foo"}' | octl iaas api CreateNic
```

## Templating

A JSON document can be used as a template, with additional config using flags.

Either from stdin:

```sh
echo '{"NetId":"vpc-foo"}' | octl iaas api CreateSubnet --IpRange 10.0.1.0/24
```

Or from a file:

```sh
octl iaas api CreateSubnet --IpRange 10.0.1.0/24 --template subnet.json
```