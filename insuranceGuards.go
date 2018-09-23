/*
    an insurance company wants to send a floor plan (an n-by-m array of "rooms"), which is pre-seeded with a set of guards (denoted by -1 ), and locked rooms (denoted by -2) and unreachable rooms (denoted by -3).  the interface should return the floor plan with each "room" (a location in the x-y plane) coded with a value > 0 with the distance to the nearest guard.


thoughts:

    - insurance is a relatively undisrupted market and is likely not basing on guard distance; to innovate in the insurance market, this is likely the wrong entry point.  additionally, most museums have a number of technical components which augment guards such as motion detection, movement (of object) sensors, laser grid, etc.  most novel attacks against museums ignore guard locations by moving the guards around or distracting them to increase guard distance.  but conceptually, a motion detector is a lot like a guard that never has to take breaks, so the analysis can be useful either way.  since insurance is fairly undisrupted, the data set for this likely already exists inside of a company like barclay's or another large insurance company in london/new york (the two major global financial hubs), and is probably adjusted based on lived experience of the firms.
    - this would be a poor market entry choice for a net-new business as software isn't the defining characteristic but rather brand trust and corporate history (rightly or wrongly)
    - reminds me of Dijkstra’s and other shortest path algorithm mechanisms, which i have no prior experience in; those use existing costs, in this case, we're calculating costs
*/
// find the distance from any room to the nearest guard, considering locked rooms

package insuranceGuards

import (
	"container/list"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

const GUARD = -2
const LOCKED = -1
const UNREACHABLE = 0

/*
   basic functionality:
       create a grid (NewGrid)
       put some guards in (PlaceGuards)
       fix the guard locations (Fix)
       score the grid (Score)
*/
// maxes are 1-indexed
type grid struct {
	// [row][colum]
	dist       [][]int
	cMax, rMax int
	fixed      bool // if the guard locations can no longer change (eg: to run a Score() and get repeatable results)
}
type room struct {
	r, c int
}

// create a new grid of arbitrary size
func NewGrid(rMax, cMax int) *grid {
	var g *grid = new(grid)
	g.dist = make([][]int, rMax)
	for r := 0; r < rMax; r++ {
		g.dist[r] = make([]int, cMax)
	}
	g.cMax = cMax
	g.rMax = rMax
	return g
}

// score the floorplan
func (g *grid) Score() {
	if !g.fixed {
		log.Print("Guards and locked rooms are not fixed...")
	}
	g.visitRoomsViaGuards(g.findGuards())
}

// fix the guard locations, eg: don't allow movement
func (g *grid) Fix() {
	g.fixed = true
}

func (g *grid) PlaceGuards(rooms ...room) {
	for _, r := range rooms {
		g.placeGuard(r)
	}
}

func (g *grid) LockRooms(rooms ...room) {
	for _, r := range rooms {
		g.lockRoom(r)
	}
}

// is a room in a grid?
func (g *grid) contains(r room) bool {
	if r.r >= 0 && r.r < g.rMax && r.c >= 0 && r.c < g.cMax {
		return true
	}
	return false
}

// is a room valid
func (g *grid) validRoom(r room) {
	if !g.contains(r) {
		log.Fatal(fmt.Sprintf("room{%d, %d} is not a valid room\n", r.r, r.c))
	}
}

// mark a room as a guard
func (g *grid) placeGuard(r room) {
	if g.fixed {
		log.Fatal("guard and locked room locations are fixed.")
	}
	g.validRoom(r)
	g.dist[r.r][r.c] = GUARD
}

// is a room guarded?
func (g *grid) guarded(r room) bool {
	g.validRoom(r)
	if g.dist[r.r][r.c] == GUARD {
		return true
	}
	return false
}

//lock a room
func (g *grid) lockRoom(r room) {
	if g.fixed {
		log.Fatal("guard and locked room locations are fixed.")
	}
	g.validRoom(r)
	g.dist[r.r][r.c] = LOCKED
}

//is a room locked?
func (g *grid) locked(r room) bool {
	g.validRoom(r)
	if g.dist[r.r][r.c] == LOCKED {
		return true
	}
	return false
}

// find all the rooms marked as guards in the distance map
func (g *grid) findGuards() []room {
	var guards []room = make([]room, 0)
	for r := 0; r < g.rMax; r++ {
		for c := 0; c < g.cMax; c++ {
			if g.guarded(room{r, c}) {
				guards = append(guards, room{r, c})
			}
		}
	}
	return guards
}
func (g *grid) String() string {
	var sb strings.Builder
	for r := 0; r < g.rMax; r++ {
		for c := 0; c < g.cMax; c++ {
			if g.locked(room{r, c}) {
				sb.WriteRune(rune('L'))
			} else if g.guarded(room{r, c}) {
				sb.WriteRune(rune('G'))
			} else {
				sb.WriteString(strconv.Itoa(g.dist[r][c]))
			}
			if c != g.cMax-1 {
				sb.WriteRune(rune(' '))
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
func (r room) proposeNeighbors() []room {
	return []room{room{r.r - 1, r.c}, room{r.r, r.c - 1}, room{r.r + 1, r.c}, room{r.r, r.c + 1}}
}

// returns the neighbors (adjacent rooms) of a room, iff in the grid and not locked or a guard
func (g *grid) neighbors(r room) []room {
	var n []room = make([]room, 0)
	var c []room = r.proposeNeighbors()
	for i := range c {
		var cr room = c[i]
		if g.contains(cr) && !g.guarded(cr) && !g.locked(cr) {
			n = append(n, cr)
		}
	}
	return n
}

// for each guard, visit all the rooms and mark their distance.  track which rooms have been visited
// for any given guard, to prevent multiple-visitation
func (g *grid) visitRoomsViaGuards(guards []room) {
	for gI := range guards {
		guard := guards[gI]
		neighbors := g.neighbors(guard)
		var toVisit *list.List = list.New() // using a linked list as a queue
		for i := range neighbors {
			toVisit.PushBack(neighbors[i])
		}
		var visited map[room]bool = make(map[room]bool) // using a hashmap as a filter on the queue
		g.calcDistances(guard, toVisit, visited)
	}
}

// iterate through the queue and calculate the distance for each room, skipping visited rooms
func (g *grid) calcDistances(guard room, toVisit *list.List, visited map[room]bool) {
	e := toVisit.Front()
	for e != nil {
		var r room = e.Value.(room)
		if visited[r] {
			e = e.Next()
			continue
		}
		g.calcDistance(guard, r, visited)
		neighbors := g.neighbors(r)
		for i := range neighbors {
			if !visited[neighbors[i]] {
				toVisit.PushBack(neighbors[i])
			}
		}
		e = e.Next()
	}
}

// calculate the distance from a guard to a room by inspecting the rooms neighbors.  skip visited rooms
func (g *grid) calcDistance(guard room, r room, visited map[room]bool) {
	if abs(guard.r-r.r)+abs(guard.c-r.c) == 1 { //they are neighbors
		g.dist[r.r][r.c] = 1
	} else { // peek at the rooms neighbors and pick the lowest non-zero value > 0
		var low int = int(math.MaxInt32)
		var neighbors []room = g.neighbors(r)
		for i := range neighbors {
			n := neighbors[i]
			if g.dist[n.r][n.c] > 0 && g.dist[n.r][n.c] < low {
				low = g.dist[n.r][n.c]
			}
		}
		// ensure only to pick a new lowest if the room is vistable by a guard
		if g.dist[r.r][r.c] == 0 {
			g.dist[r.r][r.c] = low + 1
		} else if low < g.dist[r.r][r.c] {
			g.dist[r.r][r.c] = low + 1
		}
	}
	visited[r] = true
}
func abs(n int) int {
	if n < 0 {
		return n * -1
	}
	return n
}
