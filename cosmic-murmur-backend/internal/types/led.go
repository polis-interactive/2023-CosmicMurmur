package types

import "fmt"

type LedString struct {
	LedCount    int
	StringCount int
}

type LedUniverse []LedString

type LedSegment []LedUniverse

type Light struct {
	Position Point
	Pixel    int
	Color    Color
}

func (l *Light) Print() string {
	return fmt.Sprintf("[%d], (%d, %d)", l.Pixel, l.Position.X, l.Position.Y)
}
