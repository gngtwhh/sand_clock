package sand_clock

import (
	"fmt"
	"github.com/gngtwhh/gocui/cursor"
	"github.com/gngtwhh/gocui/window"
	"sync"
	"time"
)

var SandClock [][]byte

var topOfUp, bottomOfUp, bottomOfDown int
var empty = false

const (
	sand  = byte('*')
	space = byte(' ')
)

func init() {
	topOfUp = 1
	bottomOfUp = 6
	bottomOfDown = 13
}

func Drop(stop <-chan struct{}, done *sync.WaitGroup) {
	tick := time.NewTicker(100 * time.Millisecond)
	curTop := topOfUp
	var curX, curY int
	dropping := false
	var dropOneSand = func() {
		if empty {
			return
		}
		if !dropping {
			downOneSand(&curTop)
			if empty {
				return
			}
			curX, curY = sandDrop[0], sandDrop[1]
			if SandClock[curX][curY] == space {
				SandClock[curX][curY] = sand
			}
			dropping = true
		}
	}
	var eraseOneSand = func() {
		SandClock[curX][curY] = space
	}

	defer done.Done()
	for {
		select {
		case <-stop:
			return
		case <-tick.C:
			if dropping {
				if curX+1 >= bottomOfDown {
					dropping = false
				}
				if SandClock[curX+1][curY] == space {
					eraseOneSand()
					SandClock[curX+1][curY] = sand
					curX++
				} else if SandClock[curX+1][curY] == sand {
					if SandClock[curX+1][curY-1] == space {
						eraseOneSand()
						SandClock[curX+1][curY-1] = sand
						curX, curY = curX+1, curY-1
					} else if SandClock[curX+1][curY+1] == space {
						eraseOneSand()
						SandClock[curX+1][curY+1] = sand
						curX, curY = curX+1, curY+1
					} else {
						dropping = false
					}
				}
			} else if !empty {
				dropOneSand()
			} else {
				return
			}
			if !dropping {
				dropOneSand()
			}
			refresh()
			time.Sleep(1 * time.Millisecond)
		}
	}
}
func refresh() {
	cursor.HideCursor()
	window.ClearScreen()
	cursor.GotoXY(0, 0)
	n := len(SandClock)
	for i := 0; i < n; i++ {
		fmt.Printf("%s\n", string(SandClock[i]))
	}
}
func downOneSand(curTop *int) {
	delX, delY := getDelPosOfTop(curTop)
	if delX != -1 && delY != -1 {
		SandClock[delX][delY] = space
	} else {
		empty = true
	}
}

func getDelPosOfTop(topPtr *int) (x, y int) {
	top := *topPtr
	for top <= bottomOfUp {
		i, j := 1, len(SandClock[top])-1
		for ; i <= j; i, j = i+1, j-1 {
			if SandClock[top][i] == sand {
				*topPtr = top
				return top, i
			} else if SandClock[top][j] == sand {
				*topPtr = top
				return top, j
			}
		}
		top++
	}
	*topPtr = top
	return -1, -1
}
