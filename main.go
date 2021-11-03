package main

import (
    "fmt"
    "image/color"
    "math"
    "math/rand"
    "syscall/js"
    "time"

    "test-webassembly/ball"
    "test-webassembly/vec2"
)

type Ball = ball.Ball

var balls [] Ball

var done chan struct{}

var width float64
var height float64

var ctx js.Value

func main() {
    worldScale := 1000

    // Init Canvas stuff.
    doc := js.Global().Get("document")
    canvasEl := doc.Call("getElementById", "mycanvas")
    width = doc.Get("body").Get("clientWidth").Float()
    height = doc.Get("body").Get("clientHeight").Float()
    canvasEl.Call("setAttribute", "width", width)
    canvasEl.Call("setAttribute", "height", height)

    ctx = canvasEl.Call("getContext", "2d")
    ctx.Call("scale", 1 / worldScale, 1 / worldScale)
    canvasEl.Set("width", width)
    canvasEl.Set("height", height)

    rand.Seed(time.Now().UnixNano())

    balls = append(balls, Ball{ Mass: 35 * 35 * math.Pi, Radius: 35, V: vec2.Vec2{ X: 13.7, Y: -5.7 }, Center: vec2.Vec2{ X: 40, Y: 40 }, Color: color.RGBA{ 0xff, 0x00, 0xff, 0xff } })
    balls = append(balls, Ball{ Mass: 35 * 35 * math.Pi, Radius: 35, V: vec2.Vec2{ X: -13.7, Y: -5.7 }, Center: vec2.Vec2{ X: 500, Y: 40 }, Color: color.RGBA{ 0xff, 0x00, 0xff, 0xff } })

    balls = append(balls, Ball{ Mass: 5 * 5 * math.Pi, Radius: 5, V: vec2.Vec2{ X: 13.7, Y: -5.7 }, Center: vec2.Vec2{ X: 40, Y: 340 }, Color: color.RGBA{ 0xff, 0x00, 0xff, 0xff } })
    balls = append(balls, Ball{ Mass: 35 * 35 * math.Pi, Radius: 35, V: vec2.Vec2{ X: -13.7, Y: -5.7 }, Center: vec2.Vec2{ X: 500, Y: 340 }, Color: color.RGBA{ 0xff, 0x00, 0xff, 0xff } })

    for i := 0; i < 100; i++ {
        r := rand.Float64() * 100
        m := math.Pi * r * r
        balls = append(balls, Ball{
            Mass: m,
            Radius: r,
            V: vec2.Vec2 {
                X: (rand.Float64() - 0.5) * 30,
                Y: (rand.Float64() - 0.5) * 30},
            Center: vec2.Vec2 {
                X: rand.Float64() * width,
                Y: rand.Float64() * height},
            Color: color.RGBA{uint8(rand.Float32() * 255), uint8(rand.Float32() * 255), uint8(rand.Float32() * 255), 0xff } })
    }

	var tmark float64
    var renderFrame js.Func
    renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        now := args[0].Float()
        tdiff := now - tmark
        doc.Call("getElementById", "fps").Set("innerHTML", fmt.Sprintf("FPS: %.01f", 1000/tdiff))
        tmark = now

        // Pool window size to handle resize
        curBodyW := doc.Get("body").Get("clientWidth").Float()
        curBodyH := doc.Get("body").Get("clientHeight").Float()
        if curBodyW != width || curBodyH != height {
            width, height = curBodyW, curBodyH
            canvasEl.Set("width", width)
            canvasEl.Set("height", height)
        }

        advance(tdiff)

        ctx.Set("fillStyle", "rgb(0,0,0)")
        ctx.Call("fillRect", 0, 0, width, height)

        for index := 0; index < len(balls); index++ {
            ball := &balls[index]
            ctx.Set("fillStyle", fmt.Sprintf("rgb(%d,%d,%d)", ball.Color.R, ball.Color.G, ball.Color.B));

            ctx.Call("beginPath")
            ctx.Call("arc",
                ball.Center.X,
                ball.Center.Y,
                ball.Radius,
                0, 2 * math.Pi, false );
            ctx.Call("fill")
            ctx.Call("stroke")
            ctx.Call("closePath")
        }

        js.Global().Call("requestAnimationFrame", renderFrame)
        return nil
    })

    js.Global().Call("requestAnimationFrame", renderFrame)

    println("done!")
    <-done
}

func advance(deltaT float64) {
    for index := 0; index < len(balls); index++ {
        ball := &balls[index]

        for j := index + 1; j < len(balls); j++ {
            ballOther := &balls[j]

            if ball.Center.Distance(ballOther.Center) < (ball.Radius + ballOther.Radius) {
                ball.Collide(ballOther)
            }
        }

        if ball.Center.X + ball.V.X > width - ball.Radius {
            ball.V.X = -ball.V.X
            ball.Center.X = width - ball.Radius
        } else if ball.Center.X + ball.V.X < ball.Radius {
            ball.V.X = -ball.V.X
            ball.Center.X = ball.Radius
        }
        if ball.Center.Y + ball.V.Y > height - ball.Radius {
            ball.V.Y = -ball.V.Y
            ball.Center.Y = height - ball.Radius
        } else if ball.Center.Y + ball.V.Y < ball.Radius {
            ball.V.Y = -ball.V.Y
            ball.Center.Y = ball.Radius
        }

        ball.Center = ball.Center.Plus(ball.V.Times(deltaT * 0.01))
  }
}
