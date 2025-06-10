#!/bin/sh

## activate virtual environment
. tools/sandbox/bin/activate

cd build
cmake ..
make

## load module
sudo insmod tester.ko

## unload module
sudo rmmod ldd.ko

## probe dmesg
sudo dmesg

## probe dmesg with clearing the history
sudo dmesg -c

## probe kernel modules
lsmod
