#!/bin/sh

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

python ./platform/scripts/platform-setup.py

rm -rf platform
rm -rf app
rm -rf build/components

rm platform-setup.sh
rm -rf CMakeLists.txt
