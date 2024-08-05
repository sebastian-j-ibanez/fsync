package protocol

import (
	"errors"
	"os"
	"fsync/util"
)

const kMaxBodySize = 61440

type Packet struct {
	OrderNum uint64
	BodySize uint32
	Body []byte
}

func PacketizeFile(path string, name string) ([]Packet, error) {
	// Open file
	fullPath := path + "/" + name
	file, err := os.Open(fullPath)
	util.CheckError(err)
	defer file.Close()

	// Get file size from stats
	fstats, err := file.Stat()
	util.CheckError(err)
	fsize := fstats.Size()
	if fsize <= 0 {
		return nil, errors.New("cannot packetize empty file")
	}
	
	// Calculate number of packets from file size
	pktNum := int(1)
	if fsize >= kMaxBodySize {
		pktNum = int(fsize) / int(kMaxBodySize)
	}
	if (fsize % kMaxBodySize) > 0 {
		pktNum++
	}

	// Create packets
	var packets []Packet
	offset := int64(0)
	for i := range pktNum {
		var tempBody []byte
		bytesRead, err := file.ReadAt(tempBody, offset)
		util.CheckError(err)
		packets = append(packets, Packet{
			uint64(i),
			uint32(bytesRead),
			tempBody,
		})
	}
	
	return packets, nil
}
