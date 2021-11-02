package main

import (
    "image/color"
    "math/rand"
    "syscall/js"
    "time"

    "test-webassembly/vec2"

    "github.com/llgcode/draw2d/draw2dimg"
    "github.com/llgcode/draw2d/draw2dkit"
    "github.com/markfarnan/go-canvas/canvas"
)

type Ball struct {
    color color.RGBA
    radius float64
    center vec2.Vec2
    v vec2.Vec2
}

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

    for i := 0; i < 10; i++ {
        ball := Ball{ radius: 35, v: vec2.Vec2{ X: 13.7, Y: -13.7 }, center: vec2.Vec2{ X: 40, Y: 40 }, color: color.RGBA{ 0xff, 0x00, 0xff, 0xff } }
        balls = append(balls, ball)
        balls = append(balls, Ball{
            radius: rand.Float64() * 50,
            v: vec2.Vec2 {
                X: (rand.Float64() - 0.5) * 30,
                Y: (rand.Float64() - 0.5) * 30},
            center: vec2.Vec2 {
                X: rand.Float64() * width,
                Y: rand.Float64() * height},
            color: color.RGBA{uint8(rand.Float32() * 255), uint8(rand.Float32() * 255), uint8(rand.Float32() * 255), uint8(rand.Float32() * 255)} })
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
    var backGroundColor color.RGBA = color.RGBA{0x00, 0x00, 0x00, 0xff}
    context.SetFillColor(backGroundColor)
    context.Clear()

    for index := 0; index < len(balls); index++ {
        ball := &balls[index]

        if ball.center.X + ball.v.X > width - ball.radius || ball.center.X + ball.v.X < ball.radius {
            ball.v.X = -ball.v.X
        }
        if ball.center.Y + ball.v.Y > height - ball.radius || ball.center.Y + ball.v.Y < ball.radius {
            ball.v.Y = -ball.v.Y
        }

        ball.center.X += ball.v.X
        ball.center.Y += ball.v.Y

        context.SetFillColor(ball.color)
        context.SetStrokeColor(ball.color)

        context.BeginPath()
        // context.ArcTo(ball.x, ball.y, ball.radius, ball.radius, 0, 6.12)
        draw2dkit.Circle(context, ball.center.X, ball.center.Y, ball.radius)
        context.FillStroke()
  }

  context.Close()

  return true
}
