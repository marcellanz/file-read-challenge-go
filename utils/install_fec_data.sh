#!/usr/bin/env bash

curl -L -O https://www.fec.gov/files/bulk-downloads/2018/indiv18.zip &&
unzip -p indiv18 itcont.txt | head -n 18245416 > itcont.txt &&
rm indiv18.zip &&
mkdir -p ./data/indiv18 &&
mv itcont.txt ./data/indiv18/
