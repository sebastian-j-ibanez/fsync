package protocol

import (
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/cloudflare/circl/hpke"
	"github.com/sebastian-j-ibanez/fsync/status"
)

const (
	Port        = "2000"
	MaxBodySize = 61440
)

type SocketHandler struct {
	Conn   net.Conn
	Enc    *gob.Encoder
	Dec    *gob.Decoder
	Opener hpke.Opener
	Sealer hpke.Sealer
}

// Initialize socket handler with connection
func NewSocketHandler(conn net.Conn, listenFlag bool) (SocketHandler, error) {
	var s SocketHandler

	if conn == nil {
		return SocketHandler{}, errors.New("connection is nil")
	}

	// Initialize socket connection
	s.Conn = conn
	s.Enc = gob.NewEncoder(conn)
	s.Dec = gob.NewDecoder(conn)

	if listenFlag {
		opener, sealer, err := s.setupServerEncryption()
		if err != nil {
			return SocketHandler{}, err
		}
		s.Opener = opener
		s.Sealer = sealer
	} else {
		opener, sealer, err := s.setupClientEncryption()
		if err != nil {
			return SocketHandler{}, err
		}
		s.Opener = opener
		s.Sealer = sealer
	}

	return s, nil
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

	fmt.Printf("Sending %s\n", file.Name())

	if s.Enc == nil {
		return errors.New("socket encoder uninitialized")
	}

	// Calculate and send file size
	fileSize := fileStat.Size()
	var fileSizePkt Packet
	err = fileSizePkt.SerializeToBody(fileSize, Int64)
	if err != nil {
		return err
	}

	err = s.SendEncryptedPacket(fileSizePkt)
	if err != nil {
		return err
	}

	// Calculate and send number of packets
	pktNum := CalculatePktNum(fileSize)
	var pktNumPkt Packet
	err = pktNumPkt.SerializeToBody(pktNum, Int64)
	if err != nil {
		return err
	}

	err = s.SendEncryptedPacket(pktNumPkt)
	if err != nil {
		return err
	}

	// Write incoming packets to file
	progress := status.Progress{
		TimeElapsed:    time.Now().Unix(),
		Percentage:     0,
		TotalFileBytes: fileSize,
		BytesReceived:  0,
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
		err = s.SendEncryptedPacket(tempPkt)
		if err != nil {
			return err
		}

		progress.BytesReceived += int64(bytesRead)
		progress.DisplayProgress()
	}

	fmt.Println()
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
	err = s.ReceiveEncryptedData(&fileSize, Int64)
	if err != nil {
		return err
	}

	// Get number of incoming packets
	var totalPackets int64
	err = s.ReceiveEncryptedData(&totalPackets, Int64)
	if err != nil {
		return err
	}

	fmt.Printf("Downloading %s\n", file.Name())

	// Write incoming packets to file
	progress := status.Progress{
		TimeElapsed:    time.Now().Unix(),
		Percentage:     0,
		TotalFileBytes: fileSize,
		BytesReceived:  0,
	}

	for range totalPackets {
		var tempPkt Packet
		err = s.ReceiveEncryptedPacket(&tempPkt)
		if err != nil {
			return err
		}
		bytesWritten, err := file.WriteAt(tempPkt.Body, progress.BytesReceived)
		if err != nil {
			return err
		}
		progress.BytesReceived += int64(bytesWritten)
		progress.DisplayProgress()
	}

	fmt.Println()
	return nil
}

// Send generic data over socket
func (s *SocketHandler) SendEncryptedPacket(pkt Packet) error {
	if s.Enc == nil {
		return errors.New("socket encoder uninitialized")
	}

	ct, err := s.Sealer.Seal(pkt.Body, nil)
	if err != nil {
		return nil
	}

	pkt.Body = ct

	err = s.Enc.Encode(pkt)
	if err != nil {
		return err
	}

	return nil
}

// Receive encrypted packet from socket, write to pkt
func (s *SocketHandler) ReceiveEncryptedPacket(pkt *Packet) error {
	if s.Dec == nil {
		return errors.New("socket decoder uninitialized")
	}

	err := s.Dec.Decode(&pkt)
	if err != nil {
		return err
	}

	pt, err := s.Opener.Open(pkt.Body, nil)
	if err != nil {
		return err
	}

	pkt.Body = pt

	return nil
}

// Receive encrypted data and deserialize
func (s *SocketHandler) ReceiveEncryptedData(data interface{}, pktType PacketType) error {
	var pkt Packet
	err := s.ReceiveEncryptedPacket(&pkt)
	if err != nil {
		return err
	}

	if pkt.Type != pktType {
		msg := fmt.Sprintf("packet type mismatch: expected %d, received %d", pkt.Type, pktType)
		return errors.New(msg)
	}

	err = pkt.DeserializeBody(data)
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
