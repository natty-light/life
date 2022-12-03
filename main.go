package main

import (
	"image/color"
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dgl"
)

var (
	// global rotation
	width, height    int
	redraw           = true
	gc               *draw2dgl.GraphicContext
	cellWidth        = 40
	placeMode        = true
	board            [20][20]*Cell
	cursorX, cursorY int
)

type Cell struct {
	alive      bool
	shouldLive bool
}

func reshape(window *glfw.Window, w, h int) {
	gl.ClearColor(1, 1, 1, 1)
	/* Establish viewing area to cover entire window. */
	gl.Viewport(0, 0, int32(w), int32(h))
	/* PROJECTION Matrix mode. */
	gl.MatrixMode(gl.PROJECTION)
	/* Reset project matrix. */
	gl.LoadIdentity()
	/* Map abstract coords directly to window coords. */
	gl.Ortho(0, float64(w), 0, float64(h), -1, 1)
	/* Invert Y axis so increasing Y goes down. */
	gl.Scalef(1, -1, 1)
	/* Shift origin up to upper-left corner. */
	gl.Translatef(0, float32(-h), 0)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.DEPTH_TEST)
	width, height = w, h
	redraw = true
	gc = draw2dgl.NewGraphicContext(width, height)
}

func display(gc draw2d.GraphicContext) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(0, 0, 0, 0)
	gl.LineWidth(2)

	drawBoard()
	if placeMode {
		drawCell(cursorX, cursorY, true)
	}
	if !placeMode {
		prepareNextBoard()
		updateGameState()
	}
	gl.Flush()

}

func init() {
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	width, height = 800, 800
	window, err := glfw.CreateWindow(width, height, "Life", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	window.SetSizeCallback(reshape)
	window.SetKeyCallback(onKey)
	window.SetCharCallback(onChar)

	glfw.SwapInterval(1)

	err = gl.Init()
	if err != nil {
		panic(err)
	}

	reshape(window, width, height)

	createCells()
	cursorX, cursorY = 0, 0

	for !window.ShouldClose() {
		if redraw {
			display(gc)
			window.SwapBuffers()
		}
		glfw.PollEvents()
	}
}

func onChar(w *glfw.Window, char rune) {
	log.Println(char)
}

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch {
	case key == glfw.KeyEscape && action == glfw.Press,
		key == glfw.KeyQ && action == glfw.Press:
		w.SetShouldClose(true)
	case key == glfw.KeyUp && action == glfw.Press && placeMode:
		if cursorY >= 1 {
			cursorY -= 1
		}
	case key == glfw.KeyDown && action == glfw.Press && placeMode:
		if cursorY <= width/cellWidth-2 {
			cursorY += 1
		}
	case key == glfw.KeyRight && action == glfw.Press && placeMode:
		if cursorX <= width/cellWidth-2 {
			cursorX += 1
		}
	case key == glfw.KeyLeft && action == glfw.Press && placeMode:
		if cursorX >= 1 {
			cursorX -= 1
		}
	case key == glfw.KeySpace && action == glfw.Press && placeMode:
		board[cursorX][cursorY].alive = !board[cursorX][cursorY].alive
	case key == glfw.KeyZ && action == glfw.Press:
		placeMode = !placeMode
	}
}

func drawCell(x int, y int, isCursor bool) {
	xPos, yPos := float64(x*cellWidth), float64(y*cellWidth)
	gc.MoveTo(xPos, yPos)
	gc.LineTo(xPos+float64(cellWidth), yPos)
	gc.LineTo(xPos+float64(cellWidth), yPos+float64(cellWidth))
	gc.LineTo(xPos, yPos+float64(cellWidth))
	gc.LineTo(xPos, yPos)
	gc.Close()
	gc.SetStrokeColor(color.NRGBA{0xFF, 0xFF, 0xFF, 0x80})
	if isCursor {
		gc.SetFillColor(color.NRGBA{0x80, 0x80, 0xFF, 0xFF})
	} else if board[x][y].alive {
		gc.SetFillColor(color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF})
	} else {
		gc.SetFillColor(color.NRGBA{0x00, 0x00, 0x00, 0xFF})
	}
	gc.FillStroke()
}

func createCells() {
	for x := 0; x < len(board); x++ {
		for y := 0; y < len(board[x]); y++ {
			board[x][y] = &Cell{false, false}
		}
	}
}

func drawBoard() {
	for x := 0; x < len(board); x++ {
		for y := 0; y < len(board[x]); y++ {
			drawCell(x, y, false)
		}
	}
}

func prepareNextBoard() {
	for x := 0; x < len(board); x++ {
		for y := 0; y < len(board[x]); y++ {
			applyRules(x, y)
		}
	}
}

func applyRules(xIndex int, yIndex int) {
	var neighborCount int = 0
	for x := xIndex - 1; x < xIndex+1; x++ {
		for y := yIndex - 1; y < yIndex+1; y++ {
			if x >= 0 && x <= width/cellWidth && y >= 0 && y <= width/cellWidth {
				neighborCount++
			}
		}
	}

	if board[xIndex][yIndex].alive && neighborCount <= 2 {
		board[xIndex][yIndex].shouldLive = false
	} else if board[xIndex][yIndex].alive && neighborCount <= 3 {
		board[xIndex][yIndex].shouldLive = true
	} else if board[xIndex][yIndex].alive && neighborCount > 3 {
		board[xIndex][yIndex].shouldLive = false
	} else if !board[xIndex][yIndex].alive && neighborCount == 3 {
		board[xIndex][yIndex].shouldLive = true
	}
}

func updateGameState() {
	for x := 0; x < len(board); x++ {
		for y := 0; y < len(board[x]); y++ {
			board[x][y].alive = board[x][y].shouldLive
		}
	}
}
