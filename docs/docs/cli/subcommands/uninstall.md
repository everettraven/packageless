---
id: subcommand-uninstall
title: uninstall
---

## Usage
```
packageless uninstall [PACKAGE]
```

Packages follow a particular format. If you specify just the package that you want uninstalled, the latest version of the package that **packageless** has will be uninstalled.

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
These examples do NOT reflect packages that can be used by **packageless** and is just for demonstration purposes
:::
Uninstalling the latest version of python:
```
packageless uninstall python
```
OR
```
packageless uninstall python:latest
```

Uninstalling python 3.7:
```
packageless uninstall python:3.7
```