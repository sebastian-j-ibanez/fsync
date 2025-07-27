# July 26 2025
# Sebastian J Ibanez
# This script builds and copies the fsync binary to ~/.local/bin

echo "[INFO] Building fsync binary..."
go build

echo "[INFO] Copying fsync binary to ~/.bin..."
cp .\fsync.exe ~\.bin

echo "[INFO] Setup complete!"
