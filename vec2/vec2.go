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

func (this Vec2) Mag() float64 {
    m := math.Sqrt(this.X * this.X + this.Y * this.Y)
    return m
}

func (this Vec2) Normalize() Vec2 {
    m := this.Mag()
    this.X /= m
    this.Y /= m
    return this
}

func (this Vec2) Times(scalar float64) Vec2{
    this.X *= scalar
    this.Y *= scalar
    return this
}

//
//   divided_by( scalar ) {
//     this.x /= scalar;
//     this.y /= scalar;
//     return this;
//   }
//
//
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

//
//   toString() {
//     return "(" + this.x + ", " + this.y + ")";
//   }
//
//   toStringVerbose() {
//     return "vec2 x: " + this.x + ", y: " + this.y;
//   }
//
// }
