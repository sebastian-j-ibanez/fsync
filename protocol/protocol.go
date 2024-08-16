package protocol

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

type SocketHandler struct {
	Peers []Peer
	Con net.Conn
}

type Packet struct {
	OrderNum int
	Body []byte
}

const MaxBodySize = 61440

func (s SocketHandler) SendPackets(packets []Packet) error {
	enc := gob.NewEncoder(s.Con)
	pktNum := len(packets)
	err := enc.Encode(pktNum)
	if err != nil {
		return err
	}
	for _, packet := range packets {
		err := enc.Encode(packet)
		if err != nil {
			return err
		}
	}
	
	return nil
}

func (s SocketHandler) SendPacket(packet Packet) error {
	enc := gob.NewEncoder(s.Con)
	err := enc.Encode(packet)
	return err
}

func (s SocketHandler) ReceivePackets() ([]Packet, error) {
	var packets []Packet	
	var pktNum int
	dec := gob.NewDecoder(s.Con)
	err := dec.Decode(&pktNum)
	if err != nil {
		return nil, err
	}

	for range pktNum {
		var tmpPacket Packet
		err = dec.Decode(&tmpPacket)
		if err != nil {
			return nil, err
		}
		packets = append(packets, tmpPacket)
		time.Sleep(10 * time.Millisecond)
	}

	pktLen := len(packets)
	if pktLen != pktNum {
		fmt.Printf("expected %d packets, received %d", pktLen, pktNum)
	}
	
	return packets, nil
}

func GetPacketData(packets []Packet) []byte {
	var data []byte
	for _, packet := range packets {
		data = append(data, packet.Body...)
	}
	return data
}
