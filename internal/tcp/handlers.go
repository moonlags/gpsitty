package tcp

type PosititioningPacket struct {
	Latitude  float32
	Longitude float32
	Speed     uint8
	Heading   uint16
	Timestamp int64
}
