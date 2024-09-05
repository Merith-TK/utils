#!/bin/bash

TEMPDIR=$(mktemp -d)
DEPSNOTFOUND=0
# Check if git is installed
if ! command -v git &> /dev/null; then
    echo "Git is not installed. Please install git."
    # Add installation instructions based on the user's system
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "On Linux, you can install git using the package manager."
        echo "For example, on Ubuntu, run: sudo apt-get install git"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "On macOS, you can install git using Homebrew."
        echo "Run: brew install git"
    else
        echo "Please install git manually from https://git-scm.com/downloads"
    fi
    DEPSNOTFOUND=$((DEPSNOTFOUND + 1))
fi
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install go."
    # Add installation instructions based on the user's system
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "On Linux, you can install go using the package manager."
        echo "For example, on Ubuntu, run: sudo apt-get install golang"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "On macOS, you can install go using Homebrew."
        echo "Run: brew install go"
    else
        echo "Please install go manually from https://golang.org/dl/"
    fi
    DEPSNOTFOUND=$((DEPSNOTFOUND + 1))
fi
if ! command -v make &> /dev/null; then
    echo "Make is not installed. Please install make."
    # Add installation instructions based on the user's system
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "On Linux, you can install make using the package manager."
        echo "For example, on Ubuntu, run: sudo apt-get install make"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "On macOS, you can install make using Homebrew."
        echo "Run: brew install make"
    else
        echo "Please install make manually."
    fi
    DEPSNOTFOUND=$((DEPSNOTFOUND + 1))
fi

if [ $DEPSNOTFOUND -gt 0 ]; then
    exit 1
fi

# Check if go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install go."
    # Add installation instructions based on the user's system
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "On Linux, you can install go using the package manager."
        echo "For example, on Ubuntu, run: sudo apt-get install golang"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "On macOS, you can install go using Homebrew."
        echo "Run: brew install go"
    else
        echo "Please install go manually from https://golang.org/dl/"
    fi
    exit 1
fi

# Check if make is installed
if ! command -v make &> /dev/null; then
    echo "Make is not installed. Please install make."
    # Add installation instructions based on the user's system
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "On Linux, you can install make using the package manager."
        echo "For example, on Ubuntu, run: sudo apt-get install make"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "On macOS, you can install make using Homebrew."
        echo "Run: brew install make"
    else
        echo "Please install make manually."
    fi
    exit 1
fi

# Clone the repository to a temporary folder
git clone https://github.com/Merith-TK/utils $TEMPDIR

# Change directory to the cloned repository
cd $TEMPDIR

# Run make install
make install