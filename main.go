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
        fmt.Printf("mouse down mx, my: %f, %f\n", mx, my)
        return nil
    })
    defer mouseDownEvt.Release()

    touchStartEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        e := args[0]
        touch := e.Get("changedTouches").Index(0)
        mx := touch.Get("clientX").Float()
        my := touch.Get("clientY").Float()
        world.GrabBall(mx, my)
        fmt.Printf("touch start on mx, my: %f, %f\n", mx, my)
        return nil
    })
    defer touchStartEvt.Release()

    touchMoveEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        e := args[0]
        touch := e.Get("changedTouches").Index(0)
        mx := touch.Get("clientX").Float()
        my := touch.Get("clientY").Float()
        world.MoveBall(mx, my)
        fmt.Printf("touch move on mx, my: %f, %f\n", mx, my)
        return nil
    })
    defer touchMoveEvt.Release()

    touchEndEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        world.ReleaseBall()
        println("touch ended")
        return nil
    })
    defer touchEndEvt.Release()

    mouseMoveEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        e := args[0]
        mx := e.Get("clientX").Float()
        my := e.Get("clientY").Float()
        world.MoveBall(mx, my)
        return nil
    })
    defer mouseMoveEvt.Release()

    mouseUpEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        world.ReleaseBall()
        fmt.Printf("mouse up\n")
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

        var timeSteps = map[int]float64 {
            1: 0.001,
            3: 0.002,
            4: 0.004,
            5: 0.008,
            6: 0.016,
            7: 0.032,
            8: 0.064,
            9: 0.128,
            10: 0.256,
            11: 0.512,
            12: 1.024,
            13: 2.048,
            14: 4.096,
            15: 8.192,
            16: 16.384,
        }

        world.TimeScalar = timeSteps[int(fval)]
        return nil
    })
    defer speedInputEvt.Release()

    spawnRateInputEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        evt := args[0]
        fval, err := strconv.ParseFloat(evt.Get("target").Get("value").String(), 64)
        if err != nil {
            println("Invalid value", err)
            return nil
        } else {
            println("Setting spawn rate to: ", fval)
            doc.Call("getElementById", "ball-spawn-rate-value").Set("innerHTML", fmt.Sprintf("%.01f", fval))
        }

        world.BallSpawnRate = fval
        return nil
    })
    defer speedInputEvt.Release()

    doc.Call("addEventListener", "touchstart", touchStartEvt, true)
    doc.Call("addEventListener", "touchmove", touchMoveEvt, true)
    doc.Call("addEventListener", "touchend", touchEndEvt, true)
    doc.Call("addEventListener", "keyup", keyUpEvt)
    doc.Call("addEventListener", "mousedown", mouseDownEvt)
    doc.Call("addEventListener", "mousemove", mouseMoveEvt)
    doc.Call("addEventListener", "mouseup", mouseUpEvt)
    doc.Call("getElementById", "num-balls").Call("addEventListener", "input", numBallsInputEvt)
    doc.Call("getElementById", "speed").Call("addEventListener", "input", speedInputEvt)
    doc.Call("getElementById", "ball-spawn-rate").Call("addEventListener", "input", spawnRateInputEvt)

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
