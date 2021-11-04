package main

import (
    "fmt"
    "math/rand"
    "strconv"
    "syscall/js"
    "time"

    "test-webassembly/view"
    "test-webassembly/world"
)

var done chan struct{}

func main() {
    var view view.View = view.CreateCanvasContext()

    var world world.World
    world.Init(view.Width, view.Height)

    rand.Seed(time.Now().UnixNano())

    keyUpEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        e := args[0]
        println("e.Get(\"which\").Int()", e.Get("which").Int())
        return nil
    })
    defer keyUpEvt.Release()

    mouseDownEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        e := args[0]
        // Tell world that we are grabbing the ball at coordinates x,y.
        // TODO Maybe only handle events originating from the canvas element?
        // var target js.Value = e.Get("target")
        // if target != view.CanvasEl {
        //     return nil
        // }
        mx := e.Get("clientX").Float()
        my := e.Get("clientY").Float()
        world.GrabBall(mx, my)
        fmt.Printf("mx, my: %f, %f\n", mx, my)
        return nil
    })
    defer mouseDownEvt.Release()

    mouseMoveEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        e := args[0]
        mx := e.Get("clientX").Float()
        my := e.Get("clientY").Float()
        world.MoveBall(mx, my)
        fmt.Printf("mx, my: %f, %f\n", mx, my)
        return nil
    })
    defer mouseDownEvt.Release()

    mouseUpEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        world.ReleaseBall()
        fmt.Printf("Release the ball")
        return nil
    })
    defer mouseDownEvt.Release()

    doc := js.Global().Get("document")

    numBallsInputEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        evt := args[0]
        fval, err := strconv.ParseFloat(evt.Get("target").Get("value").String(), 64)
        if err != nil {
            println("Invalid value", err)
            return nil
        } else {
            println("Setting num balls to: ", fval)
            doc.Call("getElementById", "num-balls-value").Set("innerHTML", fmt.Sprintf("%.1f", fval))
        }
        world.SetNumBalls(int(fval))
        return nil
    })
    defer numBallsInputEvt.Release()

    speedInputEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        evt := args[0]
        fval, err := strconv.ParseFloat(evt.Get("target").Get("value").String(), 64)
        if err != nil {
            println("Invalid value", err)
            return nil
        } else {
            println("Setting time scalar to: ", fval)
            doc.Call("getElementById", "speed-value").Set("innerHTML", fmt.Sprintf("%.01f", fval))
        }
        world.TimeScalar = fval
        return nil
    })
    defer speedInputEvt.Release()

    doc.Call("addEventListener", "keyup", keyUpEvt)
    doc.Call("addEventListener", "mousedown", mouseDownEvt)
    doc.Call("addEventListener", "mousemove", mouseMoveEvt)
    doc.Call("addEventListener", "mouseup", mouseUpEvt)
    doc.Call("getElementById", "num-balls").Call("addEventListener", "input", numBallsInputEvt)
    doc.Call("getElementById", "speed").Call("addEventListener", "input", speedInputEvt)

    var tmark float64
    var renderFrame js.Func
    renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        now := args[0].Float()
        tdiff := now - tmark
        doc.Call("getElementById", "fps").Set("innerHTML", fmt.Sprintf("FPS: %.01f", 1000/tdiff))
        tmark = now

        world.Advance(tdiff)

        world.Draw(&view)

        js.Global().Call("requestAnimationFrame", renderFrame)

        return nil
    })

    js.Global().Call("requestAnimationFrame", renderFrame)

    println("done!")
    <-done
}
