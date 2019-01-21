#!/bin/sh

go build readfile${1}.go && ./readfile${1} ../../../indiv18/itcont.txt 2> readfile${1}.trace.pprof
