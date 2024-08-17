package protocol

import (
	"encoding/gob"
	"net"
	"os"
)

type SocketHandler struct {
	Peers []Peer
	Con net.Conn
}

type Packet struct {
	OrderNum int64
	Body []byte
}

const MaxBodySize = 61440

// Open file at path and stream file over socket connection
func (s SocketHandler) UploadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	fileStat, err := file.Stat()
	if err != nil {
		return err
	}


	// Calculate file size and send packet num
	fileSize := fileStat.Size()
	pktNum := CalculatePktNum(fileSize)
	enc := gob.NewEncoder(s.Con)
	err = enc.Encode(pktNum)
	if err != nil {
		return err
	}

	// Iterate over file, read data, send data in packet
	offset := int64(0)
	for i := range pktNum {
		// Calculate data size if uneven amount of data left 
		var dataSize int64
		if (fileSize - offset) < MaxBodySize {
			dataSize = fileSize - offset
		} else {
			dataSize = MaxBodySize
		}

		// Read data
		data := make([]byte, dataSize)
		bytesRead, err := file.ReadAt(data, offset)
		if err != nil {
			return err
		}
		offset += int64(bytesRead)

		// Create temp packet and send over socket connection
		tempPkt := Packet{
			OrderNum: i,
			Body: data,
		}
		err = enc.Encode(tempPkt)
	}
	
	return nil
}

// Save file at path
func (s SocketHandler) DownloadFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	// Get number of incoming packets
	var pktNum int64
	dec := gob.NewDecoder(s.Con)
	err = dec.Decode(&pktNum)
	if err != nil {
		return err
	}

	// Write incoming packets to file
	offset := int64(0)
	for _ = range pktNum {
		var tempPkt Packet
		err = dec.Decode(&tempPkt)
		bytesWritten, err := file.WriteAt(tempPkt.Body, offset)
		if err != nil {
			return err
		}
		offset += int64(bytesWritten)
	}
	
	return nil
}

func (s SocketHandler) SendPktNum(pktNum int64) error {
	enc := gob.NewEncoder(s.Con)
	err := enc.Encode(pktNum)
	if err != nil {
		return err
	}
	return nil
}

func (s SocketHandler) ReceivePktNum() (int64, error) {
	var pktNum int64
	dec := gob.NewDecoder(s.Con)
	err := dec.Decode(pktNum)
	return pktNum, err
}

// Calculate number of packets based on file size 
func CalculatePktNum(fileSize int64) int64 {
	pktNum := fileSize / MaxBodySize
	if rem := fileSize % MaxBodySize; rem > 0 {
		pktNum++
	}
	return pktNum
}

// Append byte slices together
func GetPacketData(packets []Packet) []byte {
	var data []byte
	for _, packet := range packets {
		data = append(data, packet.Body...)
	}
	return data
}
