#!/bin/sh

echo "-------------------------"
echo "- packageless uninstall -"
echo "-------------------------"
echo ""
echo "removing packageless binary"
rm -f ~/bin/packageless

echo "removing .packageless folder"
rm -f -r ~/.packageless

#need to add a section here to clean up the rc file for the shell

echo "removal complete"
echo ""