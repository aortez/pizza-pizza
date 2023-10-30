package vec2

import "encoding/json"
import "math"

type Vec2 struct
{
    X float64
    Y float64
}

func (this Vec2) Copy() Vec2 {
    var copy Vec2 = Vec2{
        X: this.X,
        Y: this.Y}
    return copy
}

func (this Vec2) Plus(that Vec2) Vec2 {
    this.X += that.X
    this.Y += that.Y
    return this
}

func (this Vec2) Minus(that Vec2) Vec2 {
    this.X -= that.X
    this.Y -= that.Y
    return this
}

func (this Vec2) Distance(that Vec2) float64 {
    dx := this.X - that.X
    dy := this.Y - that.Y
    var d float64 = math.Sqrt(dx * dx + dy * dy)
    return d
}

func (vec Vec2) Mag() float64 {
    m := math.Sqrt(vec.X * vec.X + vec.Y * vec.Y)
    return m
}

func (vec Vec2) Normalize() Vec2 {
    m := vec.Mag()
    vec.X /= m
    vec.Y /= m
    return vec
}

func (vec Vec2) Times(scalar float64) Vec2{
    vec.X *= scalar
    vec.Y *= scalar
    return vec
}

func (vec Vec2) DividedBy(scalar float64) Vec2{
    vec.X /= scalar
    vec.Y /= scalar
    return vec
}

func (this Vec2) Dot(that Vec2) float64 {
    var scalarProduct float64 = this.X * that.X + this.Y * that.Y
    return scalarProduct
}

func (this Vec2) ToString() string {
    res, err := json.Marshal(this)
    if err != nil {
        return string(err.Error())
    } else {
        return string(res)
    }
}
