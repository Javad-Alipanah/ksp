package alg

import (
	"ksp/graph/model"
	"ksp/internal/pkg/alg/graph"
)

// upper left cell in the chess board is considered to have (0, 0) coordinates

func cartesianToNum(x, y, size int) int {
	return y * size + x
}

func numToCartesian(num, size int) (x int, y int) {
	x = num % size
	y = num / size
	return
}

func CalcPath(source, sink model.Point, size int) []model.Point {
	g := graph.New(size * size)
	buildGraph(-1, source.X, source.Y, size, g)
	p := g.BFS(cartesianToNum(source.X, source.Y, size), cartesianToNum(sink.X, sink.Y, size))
	if p == nil {
		return nil
	}

	var path []model.Point
	for _, v := range p {
		x, y := numToCartesian(v, size)
		point := model.Point{
			X: x,
			Y: y,
		}
		path = append(path, point)
	}
	return path
}

func buildGraph(prev, x, y, size int, g *graph.Graph) {
	if x >= size || y >= size || x < 0 || y < 0 {
		return
	}

	curr := cartesianToNum(x, y, size)
	if prev >= 0 {
		g.Connect(prev, curr)
	}

	if g.HasBeenAdded(curr) {
		return
	}

	g.AddNode(curr)

	// up right
	buildGraph(curr, x + 1, y - 2, size, g)

	// up left
	buildGraph(curr, x - 1, y - 2, size, g)

	// right down
	buildGraph(curr, x + 2, y + 1, size, g)

	// right up
	buildGraph(curr, x + 2, y - 1, size, g)

	// down right
	buildGraph(curr, x + 1, y + 2, size, g)

	// down left
	buildGraph(curr, x - 1, y + 2, size, g)

	// left down
	buildGraph(curr, x - 2, y + 1, size, g)

	// left up
	buildGraph(curr, x - 2, y - 1, size, g)
}

