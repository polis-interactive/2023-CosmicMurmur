package types

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
	MinX       int
	MaxX       int
	MinY       int
	MaxY       int
	hasBeenSet bool
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (g *Grid) TrySetMinMax(p Point) {
	if !g.hasBeenSet {
		g.MinX = p.X
		g.MinY = p.Y
		g.MaxX = p.X
		g.MaxY = p.Y
		g.hasBeenSet = true
		return
	}
	g.MinX = min(g.MinX, p.X)
	g.MinY = min(g.MinY, p.Y)
	g.MaxX = max(g.MaxX, p.X)
	g.MaxY = max(g.MaxY, p.Y)
}
