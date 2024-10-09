#!/bin/bash

# This script builds and copies the fsync binary to ~/.local/bin

echo "Building fsync binary..."
go build
echo "Copying fsync binary to ~/.local/bin..."
cp fsync ~/.local/bin
echo "Setup complete!"