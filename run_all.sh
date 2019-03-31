#!/usr/bin/env bash

for i in $(seq 0 9)
do
    ./bin/rev${i} ${1}
    ./bin/rev${i} ${1}
    ./bin/rev${i} ${1}
done