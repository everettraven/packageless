---
id: subcommand-upgrade
title: upgrade
---

## Usage
```
packageless upgrade [OPTIONAL: PIM]
```

This subcommand will upgrade the pim with the current pim information in the pim list as long as the pim is already installed. If a pim is not specified it will upgrade all installed packages.

Packages follow a particular format. If you specify just the pim that you want upgraded, the latest version of the pim that **packageless** has will be upgraded.

You can also specify a particular version by following this format:
```
pim:version
```

To manually specify that you want the latest version you can use:
```
pim:latest
```
however, **packageless** defaults to getting the latest version if one is not specified

## Examples
:::note
These examples do NOT reflect packages that can be used by **packageless** and is just for demonstration purposes
:::
upgrading the latest version of python:
```
packageless upgrade python
```
OR
```
packageless upgrade python:latest
```

upgrading python 3.7:
```
packageless upgrade python:3.7
```