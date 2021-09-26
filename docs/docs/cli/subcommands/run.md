---
id: subcommand-run
title: run
---

## Usage
```
packageless run [PACKAGE]
```

When using this subcommand, **packageless** will run the package that is specified as long as it is installed. If the package is not installed the command will exit with text stating that the package specified is not installed. If you installed a specific version of a package, you will need to use the same syntax for the package for this command as well.

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