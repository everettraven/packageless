![Build Status](https://github.com/everettraven/packageless/workflows/build/badge.svg)
![Go Version](https://img.shields.io/github/go-mod/go-version/everettraven/packageless.svg)

# Overview
**packageless** is a packageless package manager. Now, what does that mean?

**packageless** utilizes containers to isolate the packages that you would normally install with a package manager into their own environments. Due to this, all package dependencies are installed within the container and don't interact with other packages on your system. This prevents any dependency issues where one package might need a specific version of a dependency and another package needs a different version of that dependency.

## How does packageless work?

**packageless** focuses on "installing" the packages by pulling a container image and modifying the directory **packageless** is installed in to set up any volumes the container may need. In essence, **packageless** installs the package without going through the process of installing the package on the host machine, making installation much easier!

**packageless** also sets aliases whenever the package is installed so it functions exactly how you are used to the package functioning.

# Development
**packageless** is still in the very early stages of development and functionality is likely to change drastically

# Contributing
Contributing guidelines are still in development

# Documentation
Documentation is still in development
