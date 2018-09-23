package insuranceGuards

import "testing"

const (
	small  = 3
	medium = 5
	large  = 10
	xlarge = 100
)

func TestSmall(t *testing.T) {
	g := NewGrid(small, small)
	g.PlaceGuards(room{1, 1})
	g.Fix()
	g.Score()
	c := [][]int{
		[]int{2, 1, 2},
		[]int{1, GUARD, 1},
		[]int{2, 1, 2},
	}
	if !g.validate(c) {
		g.problem(t, c)
	}
}
func TestMedium(t *testing.T) {
	g := NewGrid(medium, medium)
	g.PlaceGuards(room{1, 1}, room{4, 4})
	g.Fix()
	g.Score()
	c := [][]int{
		[]int{2, 1, 2, 3, 4},
		[]int{1, GUARD, 1, 2, 3},
		[]int{2, 1, 2, 3, 2},
		[]int{3, 2, 3, 2, 1},
		[]int{4, 3, 2, 1, GUARD},
	}
	if !g.validate(c) {
		g.problem(t, c)
	}
}
func TestLarge(t *testing.T) {
	g := buildLarge()
	g.Fix()
	g.Score()
	c := [][]int{
		[]int{2, 1, 2, 3, 2, 3, 4, 5, 6, 7},
		[]int{1, GUARD, LOCKED, 2, 1, 2, 3, 4, 5, 6},
		[]int{2, 1, 2, LOCKED, GUARD, 1, 2, 3, 4, 5},
		[]int{3, 2, 1, GUARD, 1, 2, 3, 4, 5, 6},
		[]int{4, 3, 2, 1, LOCKED, 3, 2, 3, 4, 5},
		[]int{5, 4, 3, 2, 3, 2, 1, 2, 3, 4},
		[]int{6, 5, 4, 3, 2, 1, GUARD, 1, 2, 3},
		[]int{7, 6, 5, 4, 3, 2, 1, 2, 3, 4},
		[]int{8, 7, 6, 5, 4, 3, 2, 3, 4, 5},
		[]int{9, 8, 7, 6, 5, 4, 3, 4, 5, 6},
	}
	if !g.validate(c) {
		g.problem(t, c)
	}
}
func TestLargeUnreachable(t *testing.T) {
	g := buildLarge()
	rooms := []room{}
	for i := 0; i < large; i++ {
		rooms = append(rooms, room{8, i})
	}
	g.LockRooms(rooms...)
	g.Fix()
	g.Score()
	c := [][]int{
		[]int{2, 1, 2, 3, 2, 3, 4, 5, 6, 7},
		[]int{1, GUARD, LOCKED, 2, 1, 2, 3, 4, 5, 6},
		[]int{2, 1, 2, LOCKED, GUARD, 1, 2, 3, 4, 5},
		[]int{3, 2, 1, GUARD, 1, 2, 3, 4, 5, 6},
		[]int{4, 3, 2, 1, LOCKED, 3, 2, 3, 4, 5},
		[]int{5, 4, 3, 2, 3, 2, 1, 2, 3, 4},
		[]int{6, 5, 4, 3, 2, 1, GUARD, 1, 2, 3},
		[]int{7, 6, 5, 4, 3, 2, 1, 2, 3, 4},
		[]int{LOCKED, LOCKED, LOCKED, LOCKED, LOCKED, LOCKED, LOCKED, LOCKED, LOCKED, LOCKED},
		[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	if !g.validate(c) {
		g.problem(t, c)
	}
}

func buildLarge() *grid {
	g := NewGrid(large, large)
	g.PlaceGuards(room{3, 3}, room{2, 4}, room{1, 1}, room{6, 6})
	g.LockRooms(room{1, 2}, room{4, 4}, room{2, 3})
	return g
}
func (g *grid) validate(validDist [][]int) bool {
	if g.rMax != len(validDist) {
		return false
	} else if g.cMax != len(validDist[0]) {
		return false
	}
	for r := 0; r < g.rMax; r++ {
		for c := 0; c < g.cMax; c++ {
			if g.dist[r][c] != validDist[r][c] {
				return false
			}
		}
	}
	return true
}
func (g *grid) problem(t *testing.T, c [][]int) {
	t.Error("grid invalid; grid:")
	t.Error(g.dist)
	t.Error("valid grid:")
	t.Error(c)
}
