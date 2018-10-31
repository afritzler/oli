# oli - OpenStack LBaaS Imploader
CLI to remove OpenStack LBaaS instances together with all dependent objects

## Installation
```
go get github.com/afritzler/oli
```

## Usage

To use `oli` you need to `source` your `XYZ-openrc.sh` file to load the OpenStack credentials.

```
Usage:
  oli [command]

Available Commands:
  delete      Delete a LoadBalancer + everything attached
  help        Help about any command
  list        List everything LBaaS specific in your tenant
```

## Commands

### list
```
Usage:
  oli list [flags]

Flags:
      --empty   Show only LoadBalancers with no Listeners and Pool.
```