package protocol

import (
	"bytes"
	"encoding/gob"
)

type PacketType int

const (
	EncryptedPacket PacketType = iota
	FileData
	FileHashes
	Int64
	Bool
)

type Packet struct {
	OrderNum int64
	Body     []byte
	Type     PacketType
}

// SerializeToBody data into packet body
func (p *Packet) SerializeToBody(data any, packetType PacketType) error {
	p.Type = packetType

	// Gob encode data
	var dataBuf bytes.Buffer
	enc := gob.NewEncoder(&dataBuf)
	err := enc.Encode(data)
	if err != nil {
		return err
	}

	p.Body = dataBuf.Bytes()
	return nil
}

// Deserialize packet body into data
func (p *Packet) DeserializeBody(data any) error {
	dec := gob.NewDecoder(bytes.NewReader(p.Body))
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	return nil
}
