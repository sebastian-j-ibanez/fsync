#!/bin/bash

# October 7 2024
# Sebastian J Ibanez
# This script builds and copies the fsync binary to ~/.local/bin

echo "[INFO] Building fsync binary..."
build_time=$( { time go build; } 2>&1 )
real_time=$( echo "$build_time" | grep real | awk '{print $2}' )

echo "[INFO] Copying fsync binary to ~/.local/bin..."
if ! [ -f ~/.local/bin ]; then
    mkdir -p ~/.local/bin
fi
cp fsync ~/.local/bin

echo "[INFO] Setup complete!"
echo "[INFO] Build time: $real_time"