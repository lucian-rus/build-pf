#!/bin/bash

# todo:
#   * handle input arguments, agnostic of order in which they appear

# run builder
build() {
    echo "> running make all"
    
    # run python pre_build script
    python tools/scripts/pre_build.py "$1"

    cd driver
    make
    cd ..

    # run python post_build script
    python tools/scripts/post_build.py "$1"
}

# run makefile clean and remove output stuff
clean() {
    echo "> running make clean"
    cd driver
    make clean
    cd ..
    rm -rf ./output

    echo "> done"
}

# delete sandbox directory
delete_sandbox() {
    echo "> removing sandbox"
    rm -rf ./tools/sandbox

}

# evaluate given arguments
if [[ "$2" == "--clean" || "$2" == "-c" ]]; then
    clean
elif [[ "$2" == "--del-sandbox" || "$2" == "-d"  ]]; then
    delete_sandbox
else
    if [ -d "./tools/sandbox" ]; then
        echo "> sandbox exists, proceeding..."
        # find a better way to do this
        echo "> activating virtual sandbox"
        source ./tools/sandbox/bin/activate
    else
        echo "> sandbox does not exist"
        python3 -m venv tools/sandbox
        echo "> sandbox created"

        echo "> activating virtual sandbox"
        source ./tools/sandbox/bin/activate

        echo "> installing requirements"
        pip install -r ./tools/resources/requirements.txt
    fi

    build $1
fi
