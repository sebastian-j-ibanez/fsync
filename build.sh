#!/bin/bash

# This script builds and copies the fsync binary to ~/.local/bin

echo "[INFO] Building fsync binary..."
go build
echo "[INFO] Copying fsync binary to ~/.local/bin..."
cp fsync ~/.local/bin
echo "[INFO] Setup complete!"