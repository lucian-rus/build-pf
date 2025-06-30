obj-m += <target_driver_obj>.o

SRC_DIR := <target_driver_dir>src
INC_DIR := <target_driver_dir>inc
KDIR := /lib/modules/$(shell uname -r)/build
PWD := $(shell pwd)

# Tell kbuild to look for source and headers in the right place 
# -> use `src` instead of `PWD` because otherwise this is not working otherwise
# https://stackoverflow.com/questions/62420558/kernel-module-built-doesnt-include-ccflags-y
ccflags-y := -I$(src)/$(INC_DIR)

# Build rule for the object file from the src directory
<target_driver_obj>-objs := $(SRC_DIR)/<target_driver_obj>.o

all:
	make -C $(KDIR) M=$(PWD) modules

clean:
	make -C $(KDIR) M=$(PWD) clean
