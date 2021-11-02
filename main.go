package main

import (
    "image/color"
    "math"
    "math/rand"
    "syscall/js"
    "time"

    . "test-webassembly/ball"
    "test-webassembly/vec2"

    "github.com/llgcode/draw2d/draw2dimg"
    "github.com/llgcode/draw2d/draw2dkit"
    "github.com/markfarnan/go-canvas/canvas"
)

var balls [] Ball

var done chan struct{}

var myCanvas *canvas.Canvas2d
var width float64
var height float64

func main() {
    // Make Canvas 90% of window size.
    myCanvas, _ = canvas.NewCanvas2d(false)
    myCanvas.Create(int(js.Global().Get("innerWidth").Float() * 0.9), int(js.Global().Get("innerHeight").Float() * 0.9))

    height = float64(myCanvas.Height())
    width = float64(myCanvas.Width())

    rand.Seed(time.Now().UnixNano())

    balls = append(balls, Ball{ Mass: 35 * 35 * math.Pi, Radius: 35, V: vec2.Vec2{ X: 13.7, Y: -5.7 }, Center: vec2.Vec2{ X: 40, Y: 40 }, Color: color.RGBA{ 0xff, 0x00, 0xff, 0xff } })
    balls = append(balls, Ball{ Mass: 35 * 35 * math.Pi, Radius: 35, V: vec2.Vec2{ X: -13.7, Y: -5.7 }, Center: vec2.Vec2{ X: 500, Y: 40 }, Color: color.RGBA{ 0xff, 0x00, 0xff, 0xff } })

    balls = append(balls, Ball{ Mass: 5 * 5 * math.Pi, Radius: 5, V: vec2.Vec2{ X: 13.7, Y: -5.7 }, Center: vec2.Vec2{ X: 40, Y: 340 }, Color: color.RGBA{ 0xff, 0x00, 0xff, 0xff } })
    balls = append(balls, Ball{ Mass: 35 * 35 * math.Pi, Radius: 35, V: vec2.Vec2{ X: -13.7, Y: -5.7 }, Center: vec2.Vec2{ X: 500, Y: 340 }, Color: color.RGBA{ 0xff, 0x00, 0xff, 0xff } })

    println("Pi: ", math.Pi)
    for i := 0; i < 30; i++ {
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

    myCanvas.Start(60, Render)

    <-done
}

// Helper function which calls the required func (in this case 'render') every time.Duration,  Call as a go-routine to prevent blocking, as this never returns
func doEvery(d time.Duration, f func(time.Time)) {
    for x := range time.Tick(d) {
        f(x)
    }
}

// Called from the 'requestAnimationFrame' function.
// It may also be called separately from a 'doEvery' function, if the user prefers drawing to be separate from the animationFrame callback.
func Render(context *draw2dimg.GraphicContext) bool {
    var outlineColor color.RGBA = color.RGBA{0x00, 0x00, 0x00, 0x88}

    var backGroundColor color.RGBA = color.RGBA{0x00, 0x00, 0x00, 0xff}
    context.SetFillColor(backGroundColor)
    context.Clear()

    for index := 0; index < len(balls); index++ {
        ball := &balls[index]

        for j := index + 1; j < len(balls); j++ {
            ballOther := &balls[j]

            if ball.Center.Distance(ballOther.Center) < (ball.Radius + ballOther.Radius) {
                println("Collided")
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

        ball.Center = ball.Center.Plus(ball.V)

        context.SetFillColor(ball.Color)
        context.SetStrokeColor(outlineColor)

        context.BeginPath()
        // context.ArcTo(ball.x, ball.y, ball.radius, ball.radius, 0, 6.12)
        draw2dkit.Circle(context, ball.Center.X, ball.Center.Y, ball.Radius)
        context.FillStroke()
  }

  context.Close()

  return true
}
