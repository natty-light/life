package main

import (
	"log"
	"math/rand"
)

type Cell struct {
	alive      bool
	shouldLive bool
}

type Game struct {
	board      [40][40]*Cell
	placeMode  bool
	mousePlace bool
	cursorX    int
	cursorY    int
	redraw     bool
}

func (g *Game) GetBoard() [40][40]*Cell {
	return g.board
}

func (g *Game) setCursor(x, y int) {
	g.cursorX, g.cursorY = x, y
}

func (g *Game) getCursor() (int, int) {
	return g.cursorX, g.cursorY
}

func (g *Game) setPlaceMode(mode bool) {
	g.placeMode = mode
}

func (g *Game) setMousePlaceMode(mode bool) {
	g.mousePlace = mode
}

func (g *Game) togglePlaceMode() {
	g.placeMode = !g.placeMode
}

func (g *Game) isMousePlaceMode() bool {
	return g.mousePlace && g.placeMode
}

func (g *Game) isCursorPlaceMode() bool {
	return g.placeMode && !g.mousePlace
}

func (g *Game) getPlaceMode() bool {
	return g.placeMode
}

func (g *Game) toggleCell() {
	g.board[g.cursorX][g.cursorY].alive = !g.board[g.cursorX][g.cursorY].alive
}

func (g *Game) toggleMousePlaceMode() {
	g.mousePlace = !g.mousePlace
}

func (g *Game) checkCell(x, y int) bool {
	if (x < 0 && x > width/cellWidth-1) && (y < 0 && y > width/cellWidth-1) {
		log.Fatal("Cell out of bounds")
	}
	return g.board[x][y].alive
}

func (g *Game) setCellFuture(x, y int, shouldLive bool) {
	g.board[x][y].shouldLive = shouldLive
}

func (g *Game) setCell(cell *Cell, x, y int) {
	g.board[x][y] = cell
}

func (g *Game) getCell(x, y int) *Cell {
	return g.board[x][y]
}

func (g *Game) ageCell(x, y int) {
	g.getCell(x, y).alive = g.getCell(x, y).shouldLive
}

func (g *Game) updateGameState() {
	for x := 0; x < boardBound; x++ {
		for y := 0; y < boardBound; y++ {
			game.ageCell(x, y)
		}
	}
}

func (g *Game) prepareNextBoard() {
	for x := 0; x < boardBound; x++ {
		for y := 0; y < boardBound; y++ {
			g.applyRules(x, y)
		}
	}
}

func (g *Game) applyRules(xIndex int, yIndex int) {
	neighborCount := 0
	for x := xIndex - 1; x <= xIndex+1; x++ {
		for y := yIndex - 1; y <= yIndex+1; y++ {
			if x >= 0 && x <= width/cellWidth-1 && y >= 0 && y <= width/cellWidth-1 && g.checkCell(x, y) {
				neighborCount++
			}
		}
	}
	if g.checkCell(xIndex, yIndex) {
		neighborCount--
	}

	if g.checkCell(xIndex, yIndex) && neighborCount < 2 {
		g.setCellFuture(xIndex, yIndex, false)
	} else if g.checkCell(xIndex, yIndex) && neighborCount >= 2 && neighborCount <= 3 {
		g.setCellFuture(xIndex, yIndex, true)
	} else if g.checkCell(xIndex, yIndex) && neighborCount > 3 {
		g.setCellFuture(xIndex, yIndex, false)
	} else if !g.checkCell(xIndex, yIndex) && neighborCount == 3 {
		g.setCellFuture(xIndex, yIndex, true)
	}
}
func (g *Game) createCells() {
	for x := 0; x < boardBound; x++ {
		for y := 0; y < boardBound; y++ {
			g.setCell(&Cell{false, false}, x, y)
		}
	}
}

func (g *Game) createRandomPattern() {
	for x := 0; x < boardBound; x++ {
		for y := 0; y < boardBound; y++ {
			g.setCell(&Cell{rand.Intn(2) == 1, false}, x, y)
		}
	}
}

func (g *Game) drawBoard() {
	for x := 0; x < boardBound; x++ {
		for y := 0; y < boardBound; y++ {
			g.drawCell(x, y, false)
		}
	}
}

func (g *Game) drawCell(x int, y int, isCursor bool) {
	xPos, yPos := float64(x*cellWidth), float64(y*cellWidth)
	gc.MoveTo(xPos, yPos)
	gc.LineTo(xPos+float64(cellWidth), yPos)
	gc.LineTo(xPos+float64(cellWidth), yPos+float64(cellWidth))
	gc.LineTo(xPos, yPos+float64(cellWidth))
	gc.LineTo(xPos, yPos)
	gc.Close()
	gc.SetStrokeColor(cellBorderColor)
	if isCursor {
		gc.SetFillColor(cursorColor)
	} else if game.checkCell(x, y) {
		gc.SetFillColor(cellColor)
	} else {
		gc.SetFillColor(deadCellColor)
	}
	gc.FillStroke()
}

func (g *Game) getRedraw() bool {
	return g.redraw
}

func (g *Game) setRedraw(redraw bool) {
	g.redraw = redraw
}

func (g *Game) toggleRedraw() {
	g.redraw = !g.redraw
}
