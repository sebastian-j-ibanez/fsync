# fsync 🔄
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/sebastian-j-ibanez/fsync)
![GitHub Issues or Pull Requests](https://img.shields.io/github/issues/sebastian-j-ibanez/fsync?logo=github&color=blue)
![GitHub Issues or Pull Requests](https://img.shields.io/github/issues-closed/sebastian-j-ibanez/fsync?style=flat&logo=github&color=blue)
[![Go Report Card](https://goreportcard.com/badge/github.com/sebastian-j-ibanez/fsync)](https://goreportcard.com/report/github.com/sebastian-j-ibanez/fsync)

**fsync** is a CLI tool designed to sync files between devices. The protocol is entirely P2P, eliminating the need for a server or cloud instance.

## Example
Sync the ~/Pictures directory between two computers on a private network (Computer B has the pictures):
### Computer A
```
cd ~/Pictures
fsync listen --scan
```
### Computer B
```
cd ~/Pictures
fsync sync --scan
```
## Use Cases
- Sync wallpapers between devices
- Share video files

## Disclaimer
**fsync** is in early alpha. There is no gurantee on software quality or stability.

Read the [LICENSE](LICENSE) for copyright and warranty details.
