package logic

import (
	"sort"
)

// TimestampedPacket :
type TimestampedPacket struct {
	timestamp int64
	payload   []byte
}

// Packets :
type Packets struct {
	packets []TimestampedPacket
}

// NewPackets :
func NewPackets() *Packets {
	return &Packets{}
}

// AddPacket :
func (p *Packets) AddPacket(timestamp int64, payload []byte) {
	packet := TimestampedPacket{timestamp: timestamp, payload: payload}
	p.packets = append(p.packets, packet)
	sort.Slice(p.packets, func(i, j int) bool {
		return p.packets[i].timestamp > p.packets[j].timestamp
	})
}

// Peek :
func (p *Packets) Peek() *TimestampedPacket {
	if len(p.packets) == 0 {
		return nil
	}
	return &p.packets[0]
}

// Pop :
func (p *Packets) Pop() *TimestampedPacket {
	if len(p.packets) == 0 {
		return nil
	}

	v := &p.packets[0]
	elementIndex := 0
	p.packets = append(p.packets[:elementIndex], p.packets[elementIndex+1:]...)
	return v
}
