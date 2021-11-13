---
id: configuration
title: Configuring packageless
---

**packageless** has a configuration file that allows for easy configuration modifications.

The configuration file can be found in a user's home directory under the `.packageless` folder
and is called `config.hcl`.

An example to the configuration file pathing: `~/.packageless/config.hcl`

## Configuration Values
**base_dir** - The base directory that packageless should download pim configuration files, and create volumes to.

**start_port** - The port to start with when running containers that need a port exposed. This functionality is currently not enabled and this value is not used.

**port_increment** - The value to increment the start port if the port is already taken. This functionality is currently not enabled and this value is not used.

**alias** - Boolean value to indicated whether or not you would like **packageless** to automatically set aliases for you when installing a pim.

**repository_host** - The repository that **packageless** will search and pull pim configurations from. This pathing should allow for retrieving the raw file contents.

**pims_config_dir** - The directory that pim configuration files should be stored in.

**pims_dir** - The directory that volumes and date for pims to run should be created in.