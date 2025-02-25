# fsync 🔄
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/sebastian-j-ibanez/fsync?color=00ADD8)
![License](https://img.shields.io/github/license/sebastian-j-ibanez/fsync.svg?color=5E5CC4)
[![Go Report Card](https://goreportcard.com/badge/github.com/sebastian-j-ibanez/fsync)](https://goreportcard.com/report/github.com/sebastian-j-ibanez/fsync)


**fsync** is a CLI tool designed to sync files between devices. The protocol is entirely P2P, eliminating the need for a server or cloud instance.

## Example
To sync pictures from computer b -> computer a (on the same network):
### Computer A
```
cd ~/Pictures
fsync listen
```
### Computer B
```
cd ~/Pictures
fsync sync
```
## Use Cases
- Sync wallpapers between devices
- Share video files

## Disclaimer
Read the [LICENSE](LICENSE) for copyright and warranty details.
