package logic

import "math/rand"

// Sender :
type Sender interface {
	SendPacket(payload []byte)
}

// Simulate :
type Simulate struct {
	packets          Packets
	sender           Sender
	packetDropChance float32
	latency          int
}

func randBetween(first int, last int) int {
	span := last - first
	diff := int(rand.Float32()*float32(span)) + first
	return diff
}

func (s *Simulate) simulateLatency() int64 {
	span := 10
	first := s.latency - span
	last := s.latency + span
	if last > 200 {
		last = s.latency
		first = last - span*2
	}
	diff := randBetween(first, last)
	s.latency += diff
	return int64(s.latency)
}

// NewSimulate :
func NewSimulate() *Simulate {
	return &Simulate{}
}

func (s *Simulate) sendPacket(packet *TimestampedPacket) {
	s.sender.SendPacket(packet.payload)
}

func chance() float32 {
	return rand.Float32()
}

// AddPacket :
func (s *Simulate) AddPacket(now int64, payload []byte) {
	c := chance()

	if c < s.packetDropChance {
		return
	}

	const shouldDuplicate = false
	if shouldDuplicate {
		s.AddPacket(now, payload)
	}
	simulatedLatency := s.simulateLatency()
	s.packets.AddPacket(now+simulatedLatency, payload)
}

// SendPackets :
func (s *Simulate) SendPackets(now int64) {
	nextPacket := s.packets.Peek()
	if nextPacket == nil {
		return
	}

	for now >= nextPacket.timestamp {
		s.sendPacket(nextPacket)
	}
}
