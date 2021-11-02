package vec2

type Vec2 struct
{
    X float64
    Y float64
}
// 
// func CreateVec2(x float64, y float64) Vec2 {
//     var v Vec2
//     v.X = x
//     v.Y = y
//     return v
// }
//

// }
// class vec2
// {
//   constructor( x, y ) {
//     this.x = x;
//     this.y = y;
//   }
//
//   copy() {
//     let c = new vec2( this.x, this.y );
//     return c;
//   }
//
//   distance( b ) {
//     let dx = this.x - b.x;
//     let dy = this.y - b.y;
//     let d = Math.sqrt( dx * dx + dy * dy );
//     return d;
//   }
//
//   plus( a ) {
//     this.x += a.x;
//     this.y += a.y;
//     return this;
//   }
//
//   minus( a ) {
//     this.x -= a.x;
//     this.y -= a.y;
//     return this;
//   }
//
//   times( scalar ) {
//     this.x *= scalar;
//     this.y *= scalar;
//     return this;
//   }
//
//   divided_by( scalar ) {
//     this.x /= scalar;
//     this.y /= scalar;
//     return this;
//   }
//
//   mag() {
//     let m = Math.sqrt( this.x * this.x + this.y * this.y );
//     return m;
//   }
//
//   dot( b ) {
//     let scalarProduct = this.x * b.x + this.y * b.y;
//     return scalarProduct;
//   }
//
//   normalize() {
//     let m = this.mag();
//     this.x /= m;
//     this.y /= m;
//     return this;
//   }
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
