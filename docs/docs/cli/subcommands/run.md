---
id: subcommand-run
title: run
---

## Usage
```
packageless run [pim]
```

When using this subcommand, **packageless** will run the pim that is specified as long as it is installed. If the pim is not installed the command will exit with text stating that the pim specified is not installed. If you installed a specific version of a pim, you will need to use the same syntax for the pim for this command as well.

## Examples
:::note
These examples do NOT reflect packages that can be used by **packageless** and is just for demonstration purposes
:::
Running the latest version of python:
```
packageless run python
```
OR
```
packageless run python:latest
```

Running python 3.7:
```
packageless run python:3.7
```