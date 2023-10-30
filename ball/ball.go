package ball

import (
    json "encoding/json"
    "image/color"
    "math/rand"
    "test-webassembly/vec2"
)

type Ball struct {
    Color color.RGBA
    Radius float64
    Center vec2.Vec2
    Health float64
    Mass float64
    V vec2.Vec2
}

func (this *Ball) Collide (that *Ball) {
    // D = Vector between centers.
    D := this.Center.Minus(that.Center)

    // Test to see if circles are in the exact same location.
    // If so, move them a small amount so they are offset.
    distance := D.Mag()
    for {
        if distance > 0.001 {
            break
        }
        println("jittering")

        // Give the other object a small random jitter.
        that.Center.X += rand.Float64() * 0.01
        that.Center.Y += rand.Float64() * 0.01
        D = this.Center.Minus(that.Center)
        distance = D.Mag()
    }

    // Normalize vector between centers.
    Dn := D.Normalize()

    // Find min translation distance to separate circles.
    T := Dn.Times(this.Radius + that.Radius - distance)

    // Compute masses.
    m1 := this.Mass
    m2 := that.Mass
    M := m1 + m2

    // Push the circles apart, proportional to their mass.
    this.Center = this.Center.Plus(T.Times(m2 / M))
    that.Center = that.Center.Minus(T.Times(m1 / M))

    // Vector tangential to the collision plane.
    Dt := vec2.Vec2{ X: Dn.Y, Y: -Dn.X }

    // Split the velocity vector of the first ball into a normal and a tangential component in respect of the collision plane.
    v1n := Dn.Times(this.V.Dot(Dn))
    v1t := Dt.Times(this.V.Dot(Dt))

    // Split the velocity vector of the second ball into a normal and a tangential component in respect of the collision plane.
    v2n := Dn.Times(that.V.Dot(Dn));
    v2t := Dt.Times(that.V.Dot(Dt));

    // Calculate new velocity vectors of the balls, the tangential component stays the same, the normal component changes.
    elastic_factor := 0.9
    dv1t := Dn.Times((m1 - m2) / (M * v1n.Mag()) + 2 * m2 / M * v2n.Mag())
    dv2t := Dn.Times((m2 - m1) / (M * v2n.Mag()) + 2 * m1 / M * v1n.Mag())

    this.V = v1t.Plus(dv1t.Times(elastic_factor))
    that.V = v2t.Minus(dv2t.Times(elastic_factor))

    this.Health -= dv1t.Mag();
    that.Health -= dv2t.Mag();
}

func (ball *Ball) Contains(point vec2.Vec2) bool {
    return ball.Center.Distance(point) < ball.Radius
}

func (ball *Ball) Interact (ballOther *Ball) {
    gravityScalar := 1.0

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

func (this Ball) ToString() string {
    res, err := json.Marshal(this)
    if err != nil {
        return string(err.Error())
    } else {
        return string(res)
    }
}
