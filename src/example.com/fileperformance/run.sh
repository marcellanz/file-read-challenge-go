#!/bin/sh

go build readfile${1}.go && ./readfile${1} ../../../indiv18/itcont.txt
