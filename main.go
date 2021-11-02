package main

import (
    "image/color"
    "math/rand"
    "syscall/js"
    "time"

    "github.com/llgcode/draw2d/draw2dimg"
    "github.com/llgcode/draw2d/draw2dkit"
    "github.com/markfarnan/go-canvas/canvas"
)

type Ball struct {
    color color.RGBA
    radius,
    x,
    y,
    xv,
    yv float64
}

var balls [] Ball

var done chan struct{}

var myCanvas *canvas.Canvas2d
var width float64
var height float64

// This specifies how long a delay between calls to 'render'.     To get Frame Rate,   1s / renderDelay
var renderDelay time.Duration = 16 * time.Millisecond

func main() {
    FrameRate := time.Second / renderDelay
    println("requested FPS:", FrameRate)
    //canvas, _ = canvas.NewCanvas2d(true)

    // Make Canvas 90% of window size.
    myCanvas, _ = canvas.NewCanvas2d(false)
    myCanvas.Create(int(js.Global().Get("innerWidth").Float() * 0.9), int(js.Global().Get("innerHeight").Float() * 0.9))

    height = float64(myCanvas.Height())
    width = float64(myCanvas.Width())

    rand.Seed(time.Now().UnixNano())

    for i := 0; i < 100; i++ {
        ball := Ball{ radius: 35, xv: 13.7, yv: -13.7, x: 40, y: 40, color: color.RGBA{0xff, 0x00, 0xff, 0xff} }
        balls = append(balls, ball)
        balls = append(balls, Ball{
            radius: rand.Float64() * 50,
            xv: (rand.Float64() - 0.5) * 30, yv: (rand.Float64() - 0.5) * 30,
            x: rand.Float64() * width, y: rand.Float64() * height,
            color: color.RGBA{uint8(rand.Float32() * 255), uint8(rand.Float32() * 255), uint8(rand.Float32() * 255), uint8(rand.Float32() * 255)} })
    }

    myCanvas.Start(60, Render)

    //go doEvery(renderDelay, Render) // Kick off the Render function as go routine as it never returns
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

    ball := balls[index]

    if ball.x + ball.xv > width - ball.radius || ball.x + ball.xv < ball.radius {
      ball.xv = -ball.xv
    }
    if ball.y + ball.yv > height-ball.radius || ball.y + ball.yv < ball.radius {
      ball.yv = -ball.yv
    }

    ball.x += ball.xv
    ball.y += ball.yv

    context.SetFillColor(ball.color)
    context.SetStrokeColor(ball.color)

    context.BeginPath()
    // context.ArcTo(ball.x, ball.y, ball.size, ball.size, 0, math.Pi*2)
    draw2dkit.Circle(context, ball.x, ball.y, ball.radius)
    context.FillStroke()

    balls[index] = ball
  }

  context.Close()

  return true
}
