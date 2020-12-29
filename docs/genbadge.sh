#!/bin/bash

total=$1
out=$2

if (( $(echo "$total 50" | awk '{print ($1 < $2)}') )) ; then
  COLOR=red
elif (( $(echo "$total 80" | awk '{print ($1 > $2)}') )); then
  COLOR=green
else
  COLOR=orange
fi

curl "https://img.shields.io/badge/coverage-$total%25-$COLOR" > $out
