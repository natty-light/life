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
	width, height int
	redraw        = true
	gc            *draw2dgl.GraphicContext
	cellWidth           = 40
	cursor        *Cell = &Cell{0, 0, true}
)

type Cell struct {
	xIndex int
	yIndex int
	alive  bool
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

func display(gc draw2d.GraphicContext, cell *Cell) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(0, 0, 0, 0)
	gl.LineWidth(2)
	drawCell(cell)
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

	glfw.SwapInterval(0)

	err = gl.Init()
	if err != nil {
		panic(err)
	}

	reshape(window, width, height)

	for !window.ShouldClose() {
		if redraw {
			display(gc, cursor)
			window.SwapBuffers()
		}
		glfw.PollEvents()
		//		time.Sleep(2 * time.Second)
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
	case key == glfw.KeyUp && action == glfw.Press:
		if cursor.yIndex >= 1 {
			cursor.yIndex -= 1
		}
	case key == glfw.KeyDown && action == glfw.Press:
		if cursor.yIndex <= width/cellWidth-2 {
			cursor.yIndex += 1
		}
	case key == glfw.KeyRight && action == glfw.Press:
		if cursor.xIndex <= width/cellWidth-2 {
			cursor.xIndex += 1
		}
	case key == glfw.KeyLeft && action == glfw.Press:
		if cursor.xIndex >= 1 {
			cursor.xIndex -= 1
		}
	}
	log.Printf(`x: %d y: %d`, cursor.xIndex, cursor.yIndex)
}

func drawCell(cell *Cell) {
	xPos, yPos := float64(cell.xIndex*cellWidth), float64(cell.yIndex*cellWidth)
	gc.MoveTo(xPos, yPos)
	gc.LineTo(xPos+float64(cellWidth), yPos)
	gc.LineTo(xPos+float64(cellWidth), yPos+float64(cellWidth))
	gc.LineTo(xPos, yPos+float64(cellWidth))
	gc.LineTo(xPos, yPos)
	gc.Close()
	gc.SetStrokeColor(color.NRGBA{0x80, 0, 0, 0x80})
	gc.SetFillColor(color.NRGBA{0x80, 0x80, 0xFF, 0xFF})
	gc.FillStroke()
}
