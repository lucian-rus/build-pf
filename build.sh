#!/bin/sh

## activate virtual environment
. tools/sandbox/bin/activate

cd build
cmake ..
make
