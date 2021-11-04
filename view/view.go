package view

import (
    "syscall/js"
)

type View struct {
    Doc js.Value
    CanvasEl js.Value
    Ctx js.Value
    Height float64
    Width float64
}

// TODO: Create functions for drawing primatives so that the world can just call those.

func CreateCanvasContext() View {
    // Init Canvas stuff.
    worldScale := 1000.0
    doc := js.Global().Get("document")
    canvasEl := doc.Call("getElementById", "mycanvas")

    width := doc.Get("body").Get("clientWidth").Float()
    height := doc.Get("body").Get("clientHeight").Float()

    // Make the canvas as large as the client.
    canvasEl.Call("setAttribute", "width", width)
    canvasEl.Call("setAttribute", "height", height)

    ctx := canvasEl.Call("getContext", "2d")
    ctx.Call("scale", 1 / worldScale, 1 / worldScale)
    canvasEl.Set("width", width)
    canvasEl.Set("height", height)

    var view View
    view.CanvasEl = canvasEl
    view.Ctx = ctx
    view.Doc = doc
    view.Height = height
    view.Width = width
    return view
}
