package types

import "log"

type Point struct {
	X int
	Y int
}

func CreatePoint(x int, y int) Point {
	return Point{
		X: x, Y: y,
	}
}

func (p *Point) AlterPoint(newX int, newY int) {
	p.X = newX
	p.Y = newY
}

func (p *Point) IsEqual(pPrime Point) bool {
	return p.X == pPrime.X && p.Y == pPrime.Y
}

type Grid struct {
	MinX int
	MaxX int
	MinY int
	MaxY int
}

func (g *Grid) PrintGrid() {
	log.Printf("Corner 1: (%d, %d), Corner 2: (%d, %d)", g.MinX, g.MinY, g.MaxX, g.MaxY)
}
