package types

import "unsafe"

type Color struct {
	R uint8
	G uint8
	B uint8
	W uint8
}

func (c *Color) ToBits() (out uint32) {
	return uint32(c.W)<<24 | uint32(c.R)<<16 | uint32(c.G)<<8 | uint32(c.B)
}

type PixelBuffer struct {
	width  int
	minX   int
	minY   int
	stride int
	buffer []Color
}

func NewPixelBuffer(width, height, minX, minY, stride int) *PixelBuffer {
	return &PixelBuffer{
		width:  width,
		minX:   minX,
		minY:   minY,
		stride: stride,
		buffer: make([]Color, width*height*stride),
	}
}

func (pb *PixelBuffer) GetUnsafePointer() unsafe.Pointer {
	return unsafe.Pointer(&pb.buffer[0])
}

func (pb *PixelBuffer) GetPixel(p *Point) Color {
	mappedX := (p.X - pb.minX) * pb.stride
	mappedY := (p.Y - pb.minY) * pb.stride
	return pb.buffer[mappedX+mappedY*pb.width]
}

func (pb *PixelBuffer) GetPixelPointer(p *Point) *Color {
	mappedX := (p.X - pb.minX) * pb.stride
	mappedY := (p.Y - pb.minY) * pb.stride
	return &pb.buffer[mappedX+mappedY*pb.width]
}

func (pb *PixelBuffer) BlackOut() {
	for i := range pb.buffer {
		pb.buffer[i] = Color{}
	}
}
