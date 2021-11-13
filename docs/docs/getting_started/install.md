---
id: install
title: Installation
---

:::note
Currently **packageless** is only able to be installed by compiling from source and manually going through the steps to install it.

There are plans to make it so that you can install **packageless** in a much easier way in the future.
:::

## Install Docker
Currently **packageless** relies on having Docker installed. In the future, **packageless** will support multiple different container runtime options.

You can install Docker by following their installation instructions at: https://docs.docker.com/get-docker/

## Building from source
### Install Go
Since we are compiling from source and **packageless** is written in Go, you will need to make sure you have Go installed.

Go installation instructions can be found at: https://golang.org/doc/install

### Install Git or Download Source
Since we currently only support compiling from source, you will need to either have Git installed or download the source code directly from Github. We recommend using Git.

To install Git follow their installation instructions at: https://git-scm.com/book/en/v2/Getting-Started-Installing-Git

If you are using Git, once it is installed you can run (in a terminal window): 
```
git clone https://github.com/everettraven/packageless.git
```
to download the source code.

If you prefer to download the source code from Github you will have to make sure that you unzip the folder contents.

### Building packageless
In a terminal window (on Windows we recommend PowerShell), navigate into the **packageless** directory and run:
```
go build
```
If you are on Windows, this command should create a file named *packageless.exe*
If you are on Unix, this command should create a file named *packageless*

### Installing packageless
Now lets ensure the proper directories are created.

On Unix, if it doesn't exist, make the necessary directories by running:
```
mkdir ~/bin
mkdir ~/.packageless
```
On Windows, if it doesn't exist, make the necessary directories by running:

Command Prompt:
```
mkdir %USERPROFILE%\\bin
mkdir %USERPROFILE%\\.packageless
```
PowerShell:
```
mkdir ~/bin
mkdir ~/.packageless
```

Now copy the necessary files to the proper directories

On Unix run:
```
cp packageless ~/bin
cp config.hcl ~/.packageless
```

On Windows run:

Command Prompt
```
copy packageless.exe %USERPROFILE%\\bin
copy config.hcl %USERPROFILE%\\.packageless
```
PowerShell:
```
Copy-Item .\packageless.exe -Destination ~/bin
Copy-Item .\config.hcl -Destination ~/.packageless
```

Now we need to set the system PATH variable to contain the directory we are storing the **packageless** files in

:::note
Make sure you replace "shell" with the shell you're using.

For example, if you are using Bash it would be:
~/.bashrc
:::

On Unix you will need to add the following line at the end of your ~/."shell"rc file:
```
export PATH=$PATH:~/bin/packageless
```
and then run:
```
source ~/."shell"rc
```

On Windows run:

Command Prompt & PowerShell:
```
setx PATH "%PATH%;%USERPROFILE%\bin\packageless
```

**packageless** should now be installed on your machine!