#!/bin/sh

shell=$(echo ${SHELL##*/})
echo "Installing packageless on shell: $shell"

echo "Creating any needed directories"
if [ ! -d "$HOME/bin" ]
then
    mkdir ~/bin
fi

if [ ! -d "$HOME/.packageless" ]
then
    mkdir ~/.packageless
    mkdir ~/.packageless/pims_config
    mkdir ~/.packageless/pims
fi

echo "Downloading the executable..."
if [ "$OSTYPE" = "darwin" ]
then
    curl -L -o ~/bin/packageless https://github.com/everettraven/packageless/releases/latest/download/packageless-macos
else
    curl -L -o ~/bin/packageless https://github.com/everettraven/packageless/releases/latest/download/packageless-linux
fi

echo "Downloading packageless configuration file"
curl -L -o ~/.packageless/config.hcl https://github.com/everettraven/packageless/releases/latest/download/config.hcl

echo "Adding packageless to PATH by adding to: ~/."$shell"rc" 
echo "export PATH=\$PATH:~/bin/packageless" >> $HOME"/."$shell"rc"
echo "For changes to take effect, please restart your terminal"