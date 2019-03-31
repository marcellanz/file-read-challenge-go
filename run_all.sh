#!/usr/bin/env bash

start=0
end=9
for ((i=start; i<=end; i++))
do
    ./bin/rev${i} ${1}
    ./bin/rev${i} ${1}
    ./bin/rev${i} ${1}
done