#!/usr/bin/env bash

set -euo pipefail

# this script builds add.c into a shared object
cd "$(dirname "$0")"
gcc -fPIC -c add.c -o add.o
gcc -shared -o libadd.so add.o
