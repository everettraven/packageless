---
id: subcommand-install
title: install
---

## Usage
```
packageless install [PACKAGE]
```

Packages follow a particular format. If you specify just the package that you want installed, the latest version of the package that **packageless** has will be installed.

You can also specify a particular version by following this format:
```
package:version
```

To manually specify that you want the latest version you can use:
```
package:latest
```
however, **packageless** defaults to getting the latest version if one is not specified

## Examples
:::note
These examples do NOT reflect packages that can be installed by **packageless** and is just for demonstration purposes
:::
Installing the latest version of python:
```
packageless install python
```
OR
```
packageless install python:latest
```

Installing python 3.7:
```
packageless install python:3.7
```