package world

import (
    "fmt"
    "image/color"
    "math"
    "math/rand"
    // "strconv"
    // "syscall/js"
    // "time"

    "test-webassembly/ball"
    "test-webassembly/vec2"
    "test-webassembly/view"
)

type Ball = ball.Ball

type World struct {
    balls []Ball

    Width float64
    Height float64

    TimeScalar float64
}

func (w *World) Init () {
    w.TimeScalar = 0.01

    w.balls = []Ball{}
    w.balls = append(w.balls, Ball{ Mass: 35 * 35 * math.Pi, Radius: 35, V: vec2.Vec2{ X: 13.7, Y: -5.7 }, Center: vec2.Vec2{ X: 40, Y: 40 }, Color: color.RGBA{ 0xff, 0x00, 0xff, 0xff } })
    w.balls = append(w.balls, Ball{ Mass: 35 * 35 * math.Pi, Radius: 35, V: vec2.Vec2{ X: -13.7, Y: -5.7 }, Center: vec2.Vec2{ X: 500, Y: 40 }, Color: color.RGBA{ 0xff, 0x00, 0xff, 0xff } })

    w.balls = append(w.balls, Ball{ Mass: 5 * 5 * math.Pi, Radius: 5, V: vec2.Vec2{ X: 13.7, Y: -5.7 }, Center: vec2.Vec2{ X: 40, Y: 340 }, Color: color.RGBA{ 0xff, 0x00, 0xff, 0xff } })
    w.balls = append(w.balls, Ball{ Mass: 35 * 35 * math.Pi, Radius: 35, V: vec2.Vec2{ X: -13.7, Y: -5.7 }, Center: vec2.Vec2{ X: 500, Y: 340 }, Color: color.RGBA{ 0xff, 0x00, 0xff, 0xff } })

    for i := 0; i < 100; i++ {
        r := rand.Float64() * 100
        m := math.Pi * r * r
        w.balls = append(w.balls, Ball{
            Mass: m,
            Radius: r,
            V: vec2.Vec2 {
                X: (rand.Float64() - 0.5) * 30,
                Y: (rand.Float64() - 0.5) * 30},
            Center: vec2.Vec2 {
                X: rand.Float64() * w.Width,
                Y: rand.Float64() * w.Height},
            Color: color.RGBA{uint8(rand.Float32() * 255), uint8(rand.Float32() * 255), uint8(rand.Float32() * 255), 0xff } })
    }
}

func (w *World) Advance(deltaT float64) {
    gravityScalar := 1.0

    for index := 0; index < len(w.balls); index++ {
        ball := &w.balls[index]

        // Handle interactions with other balls.
        for j := index + 1; j < len(w.balls); j++ {
            ballOther := &w.balls[j]

            // They other collide or apply gravity to each other.
            if ball.Center.Distance(ballOther.Center) < (ball.Radius + ballOther.Radius) {
                ball.Collide(ballOther)
            } else {
                d := ball.Center.Distance(ballOther.Center);
                F := (gravityScalar * ball.Mass * ballOther.Mass) / (d * d);
                a := F / ball.Mass;
                D := (ballOther.Center.Minus(ball.Center)).Normalize();
                ball.V = ball.V.Plus(D.Times(a));
            }
        }

        // Bound off of window bounds.
        if ball.Center.X + ball.V.X > w.Width - ball.Radius {
            ball.V.X = -ball.V.X
            ball.Center.X = w.Width - ball.Radius
        } else if ball.Center.X + ball.V.X < ball.Radius {
            ball.V.X = -ball.V.X
            ball.Center.X = ball.Radius
        }
        if ball.Center.Y + ball.V.Y > w.Height - ball.Radius {
            ball.V.Y = -ball.V.Y
            ball.Center.Y = w.Height - ball.Radius
        } else if ball.Center.Y + ball.V.Y < ball.Radius {
            ball.V.Y = -ball.V.Y
            ball.Center.Y = ball.Radius
        }

        // Move the ball.
        ball.Center = ball.Center.Plus(ball.V.Times(deltaT * w.TimeScalar))
  }
}

func (world *World) Draw(view *view.View) {

    // TODO: make view bounds independent of the world bounds (zoom in and out!)

    // Poll window size in case things resized...
    // in which case just set the world size to match the
    // view size.
    curBodyW := view.Doc.Get("body").Get("clientWidth").Float()
    curBodyH := view.Doc.Get("body").Get("clientHeight").Float()
    if curBodyW != world.Width || curBodyH != world.Height {
        world.Width, world.Height = curBodyW, curBodyH
        view.CanvasEl.Set("width", world.Width)
        view.CanvasEl.Set("height", world.Height)
    }

    view.Ctx.Set("fillStyle", "rgb(0,0,0)")
    view.Ctx.Call("fillRect", 0, 0, world.Width, world.Height)

    for index := 0; index < len(world.balls); index++ {
        ball := &world.balls[index]
        view.Ctx.Set("fillStyle", fmt.Sprintf("rgb(%d,%d,%d)", ball.Color.R, ball.Color.G, ball.Color.B));

        view.Ctx.Call("beginPath")
        view.Ctx.Call("arc",
            ball.Center.X,
            ball.Center.Y,
            ball.Radius,
            0, 2 * math.Pi, false );
        view.Ctx.Call("fill")
        view.Ctx.Call("stroke")
        view.Ctx.Call("closePath")
    }
}
