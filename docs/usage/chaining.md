# Chaining commands

Commands may be chained, and attributes returned by a command can be reinjected in a subsequent command, using Go template syntax:

```sh
octl iaas api CreateNic --SubnetId subnet-foo | octl iaas api LinkNic -v --NicId {{.Nic.NicId}} --VmId i-foo --DeviceNumber 7
```