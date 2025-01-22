#!/usr/bin/env bash
#
# Install script for MDPreview. If your package manager installs plugins to a
# different location, adjust this script accordingly. Alternatively, you may
# manually perform the actions described below.

server_path="$HOME/.local/share/nvim/lazy/mdpreview/binary"

# check if directory exists
if ! [ -d "$server_path" ]; then
	echo "Failed to find MDPreview binary directory."
	exit
fi

echo "MDPreview binary directory found."
cd "$server_path"

# check for Go
if ! which go; then
	echo "Failed: Go is not isntalled or is not in PATH. Please see README on the GitHub."
	exit
fi

echo "Building server binary from source... Moving it to /usr/bin will require root privileges."
echo "Do this part manually if it makes you uncomfortable."

go mod tidy && go build . && sudo cp ./MDPreview /usr/bin

if which MDPreview; then
	echo "Success: MDPreview was found in PATH."
fi

echo "Finished."
