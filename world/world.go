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
type Vec2 = vec2.Vec2

type World struct {
    balls []Ball

    mouseDownBall Ball
    mouseDownIs bool
    mouseDownLocation Vec2

    NumDesiredBalls int

    BallSpawnRate float64

    Width float64
    Height float64

    TimeScalar float64
}

func (w *World) Init (width float64, height float64) {
    w.Width = width
    w.Height = height

    // TODO Set these params at start up by querying the controller's value
    // or by pushing the values here to the controller.
    w.TimeScalar = 1.024
    w.BallSpawnRate = 0.1

    w.balls = []Ball{}
    // w.balls = append(w.balls, Ball{ Mass: 35 * 35 * math.Pi, Radius: 35, V: vec2.Vec2{ X: 13.7, Y: -5.7 }, Center: vec2.Vec2{ X: 40, Y: 40 }, Color: color.RGBA{ 0xff, 0x00, 0xff, 0xff } })
    // w.balls = append(w.balls, Ball{ Mass: 35 * 35 * math.Pi, Radius: 35, V: vec2.Vec2{ X: -13.7, Y: -5.7 }, Center: vec2.Vec2{ X: 500, Y: 40 }, Color: color.RGBA{ 0xff, 0x00, 0xff, 0xff } })
    //
    // w.balls = append(w.balls, Ball{ Mass: 5 * 5 * math.Pi, Radius: 5, V: vec2.Vec2{ X: 13.7, Y: -5.7 }, Center: vec2.Vec2{ X: 40, Y: 340 }, Color: color.RGBA{ 0xff, 0x00, 0xff, 0xff } })
    // w.balls = append(w.balls, Ball{ Mass: 35 * 35 * math.Pi, Radius: 35, V: vec2.Vec2{ X: -13.7, Y: -5.7 }, Center: vec2.Vec2{ X: 500, Y: 340 }, Color: color.RGBA{ 0xff, 0x00, 0xff, 0xff } })

    // w.SetNumBalls(50)
}

func (w *World) Advance(deltaT float64) {
    // Handle interactions with other balls.
    for index := 0; index < len(w.balls); index++ {
        ball := &w.balls[index]

        for j := index + 1; j < len(w.balls); j++ {
            ballOther := &w.balls[j]

            ball.Interact(ballOther)
        }
    }

    // Collide all balls with any user held ball.
    if w.mouseDownIs {
        // heldBallCopy := world.mouseDownBall
        for index := 0; index < len(w.balls); index++ {
            ball := &w.balls[index]
            w.mouseDownBall.Interact(ball)
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

    // Constants? World parameters?
    numDivs := 2
    min_frag_radius := 4.0
    // EXPLODER_PARENT_VELOCITY_FACTOR := 0.5
    // EXPLODER_RADIAL_VELOCITY_SCALAR := 1
    EXPLODER_SIZE_FACTOR := 1.2

    // Explode the dead balls.
    // var newBalls []Ball
    // Collect the indices of the "to remove" in a map.
    var toRemove map[int]bool = make(map[int]bool)
    for i := 0; i < len(w.balls); i++ {
        ball := &w.balls[i]
        if ball.Health > 0 {
            continue
        }

        toRemove[i] = true

        divSize := ball.Radius / float64(numDivs);
        for y := ball.Center.Y - ball.Radius; y < ball.Center.Y + ball.Radius; y += divSize {
            for  x := ball.Center.X - ball.Radius; x < ball.Center.X + ball.Radius; x += divSize {
                // Don't add any more balls if we're at the desired level.
                if len(w.balls) >= int(float64(w.NumDesiredBalls) * 2.0 + 200) {
                    break
                }

                if ball.Center.Distance(vec2.Vec2{ X: x, Y: y }) > ball.Radius {
                    continue
                }

                r := divSize * EXPLODER_SIZE_FACTOR * ( 0.3 + rand.Float64() * 0.7 );
                if r < min_frag_radius {
                    continue
                }

                m := math.Pi * r * r
                w.balls = append(w.balls, Ball{
                    Health: m,
                    Mass: m,
                    Radius: r,
                    V: vec2.Vec2 {
                        X: (rand.Float64() - 0.5) * 0,
                        Y: (rand.Float64() - 0.5) * 0},
                    Center: vec2.Vec2 {
                        X: x,
                        Y: y},
                    Color: color.RGBA{uint8(rand.Float32() * 255), uint8(rand.Float32() * 255), uint8(rand.Float32() * 255), 0xff } })
            //
            //     let c = this.color.copy();
            //     c.randColor( 100 );
            //
            //     let new_ball = new Ball( x, y, r, c );
            //
            // let v = new_ball.center.copy().minus( this.center );
            // v.times( EXPLODER_RADIAL_VELOCITY_SCALAR );
            // v = v.plus( this.v.copy().times( EXPLODER_PARENT_VELOCITY_FACTOR ) );
            // v.times( Math.random() * ( EXPLODE_V_FACTOR ) );
            // new_ball.v = v;
            // new_ball.is_affected_by_gravity = true;
            // new_ball.is_moving = true;
            // new_ball.is_invincible = false;

            // frags.push( new_ball );
            }
        }
    }

    // Remove dead balls by creating a new list, excluding elements that we are supposed to remove.
    if (len(toRemove) > 0) {
        var newListOfBalls []Ball
        for i := 0; i < len(w.balls); i++ {
            _, removeMe := toRemove[i]
            if !removeMe {
                newListOfBalls = append(newListOfBalls, w.balls[i])
            }
        }
        w.balls = newListOfBalls
    }

    // Possibly spawn in a new ball.
    toAdd := w.NumDesiredBalls - len(w.balls)
    if toAdd > 0 && rand.Float64() > (1 - w.BallSpawnRate) {
        w.BallAdd(BallCreateRandom(w.Width, w.Height))
    }

    // Possibly move a ball due to user interaction.
    if w.mouseDownIs {
        // TODO There is probably a more consistent way to adjust the velocity to
        // properly match the frame rate.  Maybe if we kept track of the average
        // delta T we could guess what the actual velocity ought to be.
        throwScalar := 1.0

        w.mouseDownBall.V = w.mouseDownLocation.Minus(w.mouseDownBall.Center).Times(throwScalar)
        println("w.mouseDownBall.V: ", w.mouseDownBall.V.ToString())

        w.mouseDownBall.Center = w.mouseDownLocation
    }
}

func (w *World) BallAdd(b Ball) {
    w.balls = append(w.balls, b)
}

func BallCreateRandom(widthRange float64, heightRange float64) Ball {
    radius := rand.Float64() * ( (widthRange + heightRange) / 2 * 0.09)
    mass := math.Pi * radius * radius
    health := radius * radius * math.Pi
    return Ball{
        Health: health,
        Mass: mass,
        Radius: radius,
        V: vec2.Vec2 {
            X: (rand.Float64() - 0.5) * 30,
            Y: (rand.Float64() - 0.5) * 30},
        Center: vec2.Vec2 {
            X: rand.Float64() * widthRange,
            Y: rand.Float64() * heightRange},
        Color: color.RGBA{uint8(rand.Float32() * 255), uint8(rand.Float32() * 255), uint8(rand.Float32() * 255), 0xff } }
}

func (w *World) SetNumBalls (numBalls int) {
    // We'll spawn up to this many balls.
    w.NumDesiredBalls = numBalls

    // Either do nothing or remove some balls.
    if numBalls == len(w.balls) {
        return
    } else if numBalls < len(w.balls) {
        w.balls = w.balls[0:numBalls]
        return
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
        view.Draw(ball)
    }

    if world.mouseDownIs {
        view.Draw(&world.mouseDownBall)
    }
}

func (world *World) GrabBall(x float64, y float64) {
    if world.mouseDownIs {
        println("WARN: already held ball, releasing")
        world.mouseDownIs = false
    }

    location := vec2.Vec2{ X: x, Y: y }

    // Search through the list of balls.
    i := 0
    for ; i < len(world.balls); i++ {
        ball := &world.balls[i]

        // If we found one below the mouse cursor, remove it from the main
        // list of balls and keep track of it separately.
        if ball.Contains(location) {
            println("grabbing ball")
            world.mouseDownIs = true
            // TODO ASSERT if we're already holding a ball
            world.mouseDownBall = *ball
            world.mouseDownLocation = location

            // Slice the one we're holding out of the main collection.
            world.balls = append(world.balls[0 : i], world.balls[i + 1 : ]...)
            break
        }
    }

    // If we didn't grab an existing ball, then lets spawn one.
    if !world.mouseDownIs {
        println("spawning ball")
        world.mouseDownBall = BallCreateRandom(world.Width, world.Height)
        world.mouseDownBall.V = vec2.Vec2{}
        world.mouseDownBall.Center = location
        world.mouseDownLocation = location
        world.mouseDownIs = true
    }
}

func (world *World) MoveBall(x float64, y float64) {
    if !world.mouseDownIs {
        return
    }

    world.mouseDownLocation = vec2.Vec2{ X: x, Y: y }
}

func (world *World) ReleaseBall() {
    world.mouseDownIs = false
    world.balls = append(world.balls, world.mouseDownBall)
}
