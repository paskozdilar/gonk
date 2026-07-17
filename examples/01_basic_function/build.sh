#!/usr/bin/env bash

set -euo pipefail

# this script builds add.c into a shared object
cd "$(dirname "$0")"
mkdir -p ./obj ./lib
gcc -fPIC -c src/add.c -o obj/add.o
gcc -shared -o lib/libadd.so obj/add.o
