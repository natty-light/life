package main

import (
	"image/color"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dgl"
)

var (
	// global rotation
	width, height    int = 800, 800
	redraw               = true
	gc               *draw2dgl.GraphicContext
	cellWidth        = 20
	placeMode        = true
	mousePlace       = false
	board            [40][40]*Cell
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
	drawCursor()

	if !placeMode {
		prepareNextBoard()
		go updateGameState()
		time.Sleep(10 * time.Millisecond)
	}
	gl.Flush()

}

func drawCursor() {
	if placeMode && !mousePlace {
		drawCell(cursorX, cursorY, true)
	}
	if mousePlace && placeMode {
		mX, mY := glfw.GetCurrentContext().GetCursorPos()
		cursorX, cursorY = int(mX/float64(cellWidth)), int(mY/float64(cellWidth))
		drawCell(cursorX, cursorY, true)
	}
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
	window, err := glfw.CreateWindow(width, height, "Life", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	window.SetSizeCallback(reshape)
	window.SetKeyCallback(onKey)
	window.SetCharCallback(onChar)
	window.SetMouseButtonCallback(onClick)
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

// Mouse button click callback
func onClick(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if mousePlace && button == glfw.MouseButton1 && action == glfw.Press {
		board[cursorX][cursorY].alive = !board[cursorX][cursorY].alive
	}
}

// Keyboard char callback
func onChar(w *glfw.Window, char rune) {
	log.Println(char)
}

// Keyboard key callback
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
	case key == glfw.KeyM && action == glfw.Press:
		if placeMode {
			mousePlace = !mousePlace
		}
	case key == glfw.KeyC && action == glfw.Press:
		createCells()
	case key == glfw.KeyR && action == glfw.Press:
		createRandomPattern()
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

func createRandomPattern() {
	for x := 0; x < len(board); x++ {
		for y := 0; y < len(board[x]); y++ {
			board[x][y] = &Cell{rand.Intn(2) == 1, false}
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
	for x := xIndex - 1; x <= xIndex+1; x++ {
		for y := yIndex - 1; y <= yIndex+1; y++ {
			if x >= 0 && x <= width/cellWidth-1 && y >= 0 && y <= width/cellWidth-1 && board[x][y].alive {
				neighborCount++
			}
		}
	}
	if board[xIndex][yIndex].alive {
		neighborCount--
	}

	if board[xIndex][yIndex].alive && neighborCount < 2 {
		board[xIndex][yIndex].shouldLive = false
	} else if board[xIndex][yIndex].alive && neighborCount >= 2 && neighborCount <= 3 {
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
