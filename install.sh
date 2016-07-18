#!/bin/sh

go install
sudo rm -Rf /var/www/html/depense/res; sudo cp -Rf res /var/www/html/depense/
depense


