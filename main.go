package main

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
	"image/color"
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dgl"
)

var (
	// global rotation
	width, height                                              = 800, 900
	boardWidth, boardHeight                                    = width, height - 100
	boardBound                                                 = 40
	gc                                                         *draw2dgl.GraphicContext
	cellWidth                                                  = 20
	game                                                       *Game
	maxCursorX, maxCursorY                                             = boardWidth / cellWidth, boardHeight / cellWidth
	fontColumnOne, fontColumnTwo, fontSpacing, fontStartHeight float64 = 50, 400, 20, 820
)

var (
	fontColor       = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	cursorColor     = color.NRGBA{R: 0x80, G: 0x80, B: 0xFF, A: 0xFF}
	cellColor       = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	deadCellColor   = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
	cellBorderColor = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0x80}
)

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
	game.setRedraw(true)
	gc = draw2dgl.NewGraphicContext(width, height)
}

func display(gc draw2d.GraphicContext) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(0, 0, 0, 0)
	gl.LineWidth(2)

	game.drawBoard()
	drawCursor()
	drawControls()

	if !game.getPlaceMode() {
		game.prepareNextBoard()
		go game.updateGameState()
		time.Sleep(10 * time.Millisecond)
	}
	gl.Flush()
}

func drawCursor() {
	if game.isCursorPlaceMode() {
		game.drawCell(game.cursorX, game.cursorY, true)
	} else if game.isMousePlaceMode() {
		mX, mY := glfw.GetCurrentContext().GetCursorPos()
		handleCursorBoundaries(mX, mY)
		game.drawCell(game.cursorX, game.cursorY, true)
	}
}

func handleCursorBoundaries(mX, mY float64) {
	adjustedX, adjustedY := int(mX/float64(cellWidth)), int(mY/float64(cellWidth))
	if adjustedX < 0 {
		adjustedX = 0
	} else if adjustedX > maxCursorX {
		adjustedX = maxCursorX - 1
	}

	if adjustedY < 0 {
		adjustedY = 0
	} else if adjustedY > maxCursorY {
		adjustedY = maxCursorY - 1
	}

	game.setCursor(adjustedX, adjustedY)
}

func drawControls() {
	gc.SetFontData(draw2d.FontData{Name: "GoRegular", Family: draw2d.FontFamilySans, Style: draw2d.FontStyleNormal})
	gc.SetFontSize(12)
	gc.SetFillColor(fontColor)
	gc.FillStringAt("Controls", fontColumnOne, fontStartHeight)
	gc.FillStringAt("Arrow keys: Move cursor", fontColumnOne, fontStartHeight+fontSpacing)
	gc.FillStringAt("Space: Toggle cell", fontColumnTwo, fontStartHeight+fontSpacing)
	gc.FillStringAt("Z: Toggle place mode", fontColumnOne, fontStartHeight+2*fontSpacing)
	gc.FillStringAt("M: Toggle mouse placement", fontColumnTwo, fontStartHeight+2*fontSpacing)
	gc.FillStringAt("C: Clear board", fontColumnOne, fontStartHeight+3*fontSpacing)
	gc.FillStringAt("R: Random pattern", fontColumnTwo, fontStartHeight+3*fontSpacing)
}

func init() {
	runtime.LockOSThread()

	fc := FontCache{}
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	fc.Store(draw2d.FontData{Name: "GoRegular", Family: draw2d.FontFamilySans, Style: draw2d.FontStyleNormal}, font)
	draw2d.SetFontCache(fc)

	game = &Game{placeMode: true, mousePlace: false, cursorX: 0, cursorY: 0}
}

func main() {
	err := glfw.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()
	window, err := glfw.CreateWindow(width, height, "Life", nil, nil)
	if err != nil {
		log.Fatal(err)
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

	draw2d.SetFontFolder("")

	reshape(window, width, height)

	game.createCells()
	game.setCursor(0, 0)

	for !window.ShouldClose() {
		if game.getRedraw() {
			display(gc)
			window.SwapBuffers()
		}
		glfw.PollEvents()
	}
}

// Mouse button click callback
func onClick(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if game.isMousePlaceMode() && button == glfw.MouseButton1 && action == glfw.Press {
		game.toggleCell()
	}
}

// Keyboard char callback
func onChar(w *glfw.Window, char rune) {
	log.Println(char)
}

// Keyboard key callback
func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	x, y := game.getCursor()
	switch {
	case key == glfw.KeyEscape && action == glfw.Press,
		key == glfw.KeyQ && action == glfw.Press:
		w.SetShouldClose(true)
	case key == glfw.KeyUp && action == glfw.Press && game.getPlaceMode():
		if x >= 1 {
			game.setCursor(x, y-1)
		}
	case key == glfw.KeyDown && action == glfw.Press && game.getPlaceMode():
		if y <= width/cellWidth-2 {
			game.setCursor(x, y+1)
		}
	case key == glfw.KeyRight && action == glfw.Press && game.getPlaceMode():
		if x <= width/cellWidth-2 {
			game.setCursor(x+1, y)
		}
	case key == glfw.KeyLeft && action == glfw.Press && game.getPlaceMode():
		if x >= 1 {
			game.setCursor(x-1, y)
		}
	case key == glfw.KeySpace && action == glfw.Press && game.getPlaceMode():
		game.toggleCell()
	case key == glfw.KeyZ && action == glfw.Press:
		game.togglePlaceMode()
	case key == glfw.KeyM && action == glfw.Press:
		if game.getPlaceMode() {
			game.toggleMousePlaceMode()
		}
	case key == glfw.KeyC && action == glfw.Press:
		game.createCells()
	case key == glfw.KeyR && action == glfw.Press:
		game.createRandomPattern()
	}
}
