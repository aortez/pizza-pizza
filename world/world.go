package world

import (
    "fmt"
    "image/color"
    "math"
    "math/rand"

    "test-webassembly/ball"
    "test-webassembly/vec2"
    "test-webassembly/view"
)

type Ball = ball.Ball

type World struct {
    balls []Ball

    mouseDownBall *Ball

    Width float64
    Height float64

    TimeScalar float64
}

func (w *World) Init (width float64, height float64) {
    w.Width = width
    w.Height = height

    w.TimeScalar = 0.3

    w.balls = []Ball{}
    w.balls = append(w.balls, Ball{ Mass: 35 * 35 * math.Pi, Radius: 35, V: vec2.Vec2{ X: 13.7, Y: -5.7 }, Center: vec2.Vec2{ X: 40, Y: 40 }, Color: color.RGBA{ 0xff, 0x00, 0xff, 0xff } })
    w.balls = append(w.balls, Ball{ Mass: 35 * 35 * math.Pi, Radius: 35, V: vec2.Vec2{ X: -13.7, Y: -5.7 }, Center: vec2.Vec2{ X: 500, Y: 40 }, Color: color.RGBA{ 0xff, 0x00, 0xff, 0xff } })

    w.balls = append(w.balls, Ball{ Mass: 5 * 5 * math.Pi, Radius: 5, V: vec2.Vec2{ X: 13.7, Y: -5.7 }, Center: vec2.Vec2{ X: 40, Y: 340 }, Color: color.RGBA{ 0xff, 0x00, 0xff, 0xff } })
    w.balls = append(w.balls, Ball{ Mass: 35 * 35 * math.Pi, Radius: 35, V: vec2.Vec2{ X: -13.7, Y: -5.7 }, Center: vec2.Vec2{ X: 500, Y: 340 }, Color: color.RGBA{ 0xff, 0x00, 0xff, 0xff } })

    w.SetNumBalls(50)
}

func (w *World) SetNumBalls (numBalls int) {
    // Simple cases are do nothing or remove balls.
    if numBalls == len(w.balls) {
        return
    } else if numBalls < len(w.balls) {
        w.balls = w.balls[0:numBalls]
        return
    }

    // Otherwise we're adding balls.
    toAdd := numBalls - len(w.balls)
    for i := 0; i < toAdd; i++ {
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

    // Handle interactions with other balls.
    for index := 0; index < len(w.balls); index++ {
        ball := &w.balls[index]

        for j := index + 1; j < len(w.balls); j++ {
            ballOther := &w.balls[j]

            // They either collide or apply gravity to each other.
            if ball.Center.Distance(ballOther.Center) < (ball.Radius + ballOther.Radius) {
                ball.Collide(ballOther)
            } else {
                d := ball.Center.Distance(ballOther.Center);
                F := (gravityScalar * ball.Mass * ballOther.Mass) / (d * d);
                a := F / ball.Mass;
                b := F / ballOther.Mass;
                D := (ballOther.Center.Minus(ball.Center)).Normalize();
                ball.V = ball.V.Plus(D.Times(a));
                ballOther.V = ballOther.V.Minus(D.Times(b));
            }
        }
    }

    // Bounce off of window bounds.
    for index := 0; index < len(w.balls); index++ {
        ball := &w.balls[index]

        if ball.Center.X > w.Width - ball.Radius {
            ball.V.X = -ball.V.X
            ball.Center.X = w.Width - ball.Radius
        } else if ball.Center.X < ball.Radius {
            ball.V.X = -ball.V.X
            ball.Center.X = ball.Radius
        }

        if ball.Center.Y > w.Height - ball.Radius {
            ball.V.Y = -ball.V.Y
            ball.Center.Y = w.Height - ball.Radius
        } else if ball.Center.Y < ball.Radius {
            ball.V.Y = -ball.V.Y
            ball.Center.Y = ball.Radius
        }

        // Move the ball.
        timeScalar := w.TimeScalar * 0.01
        ball.Center = ball.Center.Plus(ball.V.Times(deltaT * timeScalar))
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
        fmt.Printf("Changing world.Width/Height from %f/%f to %f/%f", world.Width, world.Height, curBodyW, curBodyH)
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

func (world *World) GrabBall(x float64, y float64) {
    for i := 0; i < len(world.balls); i++ {
        ball := &world.balls[i]
        if ball.Contains(vec2.Vec2{ X: x, Y: y }) {
            println("yes we can")
            world.mouseDownBall = ball
        }
    }
}

func (world *World) MoveBall(x float64, y float64) {
    if world.mouseDownBall == nil {
        return
    }

    var ball *Ball = world.mouseDownBall

    newLocation := vec2.Vec2{ X: x, Y: y }

    // TODO There is probably a more consistent way to adjust the velocity to
    // properly match the frame rate.  Maybe if we kept track of the average
    // delta T we could guess what the actual velocity ought to be.
    throwScalar := 100
    ball.V = newLocation.Minus(ball.Center).Times(throwScalar)

    ball.Center = newLocation
}

func (world *World) ReleaseBall() {
    world.mouseDownBall = nil
}
