package protocol

import (
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/sebastian-j-ibanez/fsync/status"
)

const (
	Port        = "2000"
	MaxBodySize = 61440
)

type SocketHandler struct {
	Conn net.Conn
	Enc  *gob.Encoder
	Dec  *gob.Decoder
}

type Packet struct {
	OrderNum int64
	Body     []byte
}

// Initialize socket handler with connection
func NewSocketHandler(conn net.Conn) SocketHandler {
	var s SocketHandler
	if conn != nil {
		s.Conn = conn
		s.Enc = gob.NewEncoder(conn)
		s.Dec = gob.NewDecoder(conn)
	} else {
		s.Conn, s.Enc, s.Dec = nil, nil, nil
	}

	return s
}

// Open file at path and stream file over socket connection
func (s SocketHandler) UploadFile(path string) error {
	// Get file stats
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	fileStat, err := file.Stat()
	if err != nil {
		return err
	}

	if s.Enc == nil {
		return errors.New("socket encoder uninitialized")
	}

	// Calculate and send file size
	fileSize := fileStat.Size()
	err = s.Enc.Encode(fileSize)
	if err != nil {
		return err
	}

	// Calculate and send number of packets
	pktNum := CalculatePktNum(fileSize)
	err = s.Enc.Encode(pktNum)
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
			Body:     data,
		}
		err = s.Enc.Encode(tempPkt)
		if err != nil {
			return err
		}
	}

	return nil
}

// Save file at path
func (s *SocketHandler) DownloadFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	if s.Dec == nil {
		return errors.New("socket decoder uninitialized")
	}

	// Get file size
	var fileSize int64
	err = s.Dec.Decode(&fileSize)
	if err != nil {
		return err
	}

	// Get number of incoming packets
	var totalPackets int64
	err = s.Dec.Decode(&totalPackets)
	if err != nil {
		return err
	}

	fmt.Printf("Downloading %s\n", file.Name())

	// Write incoming packets to file
	bytesReceived := int64(0)
	for range totalPackets {
		var tempPkt Packet
		err = s.Dec.Decode(&tempPkt)
		if err != nil {
			return err
		}
		bytesWritten, err := file.WriteAt(tempPkt.Body, bytesReceived)
		if err != nil {
			return err
		}
		bytesReceived += int64(bytesWritten)
		//fmt.Printf("Bytes received: %d, total bytes: %d\n", bytesReceived, fileSize)
		status.PrintLoadingBar(bytesReceived, fileSize)
	}

	return nil
}

// Send generic data over socket
func (s *SocketHandler) SendGenericData(data any) error {
	if s.Enc == nil {
		return errors.New("socket encoder uninitialized")
	}
	err := s.Enc.Encode(data)
	if err != nil {
		return err
	}

	return nil
}

// Receive generic data from socket
func (s *SocketHandler) ReceiveGenericData(data any) error {
	if s.Dec == nil {
		return errors.New("socket decoder uninitialized")
	}
	err := s.Dec.Decode(data)
	if err != nil {
		return err
	}

	return nil
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
