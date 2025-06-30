#!/bin/sh

if [ -d "./tools/sandbox" ]; then
    echo "> sandbox exists, proceeding..."
    # find a better way to do this
    echo "> activating virtual sandbox"
    source ./tools/sandbox/bin/activate
else
    echo "> sandbox does not exist"
    mkdir tools
    python3 -m venv tools/sandbox
    echo "> sandbox created"

    echo "> activating virtual sandbox"
    source ./tools/sandbox/bin/activate
fi

python ./platform/scripts/platform-setup.py

# finish up set-up by removing unwanted files -> @todo make this configurable ??
rm -rf .git
rm -rf platform
rm platform-setup.sh
