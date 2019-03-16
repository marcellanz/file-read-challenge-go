#!/usr/bin/env bash

for i in {0..10}
do
    ./bin/rev${i} ${1} > /dev/null
    ./bin/rev${i} ${1} > /dev/null
    ./bin/rev${i} ${1} > /dev/null
done