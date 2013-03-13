package GoLife

import (
	"fmt"
	"time"
	"math/rand"
)

const (
	DX = 40
	DY = 40
	SLEEP = 1
)

type LifeGrid struct {
	grid [][]bool
	dx, dy int
}

func New(dx int, dy int) *LifeGrid {
	lg := &LifeGrid{ nil, dx, dy }

	grid := make([][]bool, dy)
	for i, _ := range grid {
		grid[i] = make([]bool, dx)
	}
	lg.grid = grid

	return lg
}

func (lg *LifeGrid) Randomize() {
	rand.Seed(time.Now().Unix())

	var i int = 0
	for y := 0; y < lg.dy; y++ {
		for x := 0; x < lg.dx; x++ {
			i = rand.Intn(2)
			if i == 0 {
				lg.grid[y][x] = false
			} else {
				lg.grid[y][x] = true
			}
		}
	}
}

func (lg LifeGrid) String(aliveChar rune, deadChar rune) string {
	var o string = ""
	for _, row := range lg.grid {
		for _, cell := range row {
			if cell {
				o += fmt.Sprintf("%s", string(aliveChar))
			} else {
				o += fmt.Sprintf("%s", string(deadChar))
			}
			o += fmt.Sprintf(" ")
		}
		o += fmt.Sprint("\n")
	}

	return o
}

func (lg *LifeGrid) Tick() {
	nc := make([][]int, len(lg.grid)) // count of neighbours
	for y, row := range lg.grid {
		nc[y] = make([]int, len(row))
		for x, _ := range row {
			nc[y][x] = countNeighbours(*lg, x, y)
		}
	}

	// update cells
	for y, row := range lg.grid {
		for x, cell := range row {
			if cell {
				switch {
				case nc[y][x] < 2:
					lg.grid[y][x] = false
				case nc[y][x] > 3:
					lg.grid[y][x] = false
				}
			} else {
				if nc[y][x] == 3 {
					lg.grid[y][x] = true
				}
			}
		}
	}
}

func countNeighbours(lg LifeGrid, x int, y int) int {
	NEIGHBOURS := [][]int {
		{-1, -1},
		{0, -1},
		{1, -1},
		{1, 0},
		{1, 1},
		{0, 1},
		{-1, 1},
		{-1, 0},
	}
	dx, dy := lg.dx, lg.dy

	n := 0
	for _, dd := range NEIGHBOURS {
		xx, yy := (x+dd[0]), (y+dd[1])
		if xx < 0 {
			xx = dx + xx 
		} else if xx >= dx {
			xx = xx % dx
		}
		if yy < 0 {
			yy = dy + yy 
		} else if yy >= dy {
			yy = yy % dy
		}

		cell := lg.grid[yy][xx]
		if cell {
			n++
		}
	}

	return n
}

func LifeTest(killChan chan int, killConfirm chan int) {

	lgrid := New(40, 40)
	lgrid.Randomize()

	for {
		select {
		case <-killChan:
			killConfirm <- 1
			break
		default:
			lgrid.Tick()

			for i := 0; i < 10; i++ {
				fmt.Print("=")
			}
			fmt.Print("\n")
			fmt.Printf("%s", lgrid.String('@', '.'))
			for i := 0; i < 10; i++ {
				fmt.Print("=")
			}
			fmt.Print("\n")

			time.Sleep(time.Duration(SLEEP) * time.Second)
		}
	}
}

