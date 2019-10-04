package geometry

import (
	"math"
	"math/rand"
)

// Vec3 is a Vector3 representation
type Vec3 struct {
	e [3]float64
}

// xyz getters

// NewVec3 create a new object of Vec3
func NewVec3(x, y, z float64) *Vec3 {
	return &Vec3{[3]float64{x, y, z}}
}

//Pv return the pointer value
func (v Vec3) Pv() *Vec3 {
	return &v
}

// X return the X value
func (v *Vec3) X() float64 {
	return v.e[0]
}

// Y return the Y value
func (v *Vec3) Y() float64 {
	return v.e[1]
}

// Z return the Z value
func (v *Vec3) Z() float64 {
	return v.e[2]
}

// R return the R value
func (v *Vec3) R() float64 {
	return v.e[0]
}

// G return the G value
func (v *Vec3) G() float64 {
	return v.e[1]
}

// B return the B value
func (v *Vec3) B() float64 {
	return v.e[2]
}

// Opposite return the opposite vector
func (v *Vec3) Opposite() *Vec3 {
	return NewVec3(-v.e[0], -v.e[1], -v.e[2])
}

// At return the value at index i
func (v *Vec3) At(i int) float64 {
	return v.e[i]
}

func (v *Vec3) SetAt(i int, value float64) {
	v.e[i] = value
}

// Plus return a new Vec3 result of the operation of this and another
func (v *Vec3) Plus(v2 *Vec3) *Vec3 {
	return NewVec3(v.e[0]+v2.e[0], v.e[1]+v2.e[1], v.e[2]+v2.e[2])
}

// PlusScalar return a new Vec3 result of the operation of this and a Scalar
func (v *Vec3) PlusScalar(t float64) *Vec3 {
	return NewVec3(v.e[0]+t, v.e[1]+t, v.e[2]+t)
}

func (v *Vec3) Minus(v2 *Vec3) *Vec3 {
	return NewVec3(v.e[0]-v2.e[0], v.e[1]-v2.e[1], v.e[2]-v2.e[2])
}

func (v *Vec3) MinusScalar(t float64) *Vec3 {
	return NewVec3(v.e[0]-t, v.e[1]-t, v.e[2]-t)
}

func (v *Vec3) Times(v2 *Vec3) *Vec3 {
	return NewVec3(v.e[0]*v2.e[0], v.e[1]*v2.e[1], v.e[2]*v2.e[2])
}

func (v *Vec3) TimesScalar(t float64) *Vec3 {
	return NewVec3(v.e[0]*t, v.e[1]*t, v.e[2]*t)
}

func (v *Vec3) By(v2 *Vec3) *Vec3 {
	return NewVec3(v.e[0]/v2.e[0], v.e[1]/v2.e[1], v.e[2]/v2.e[2])
}

func (v *Vec3) ByScalar(t float64) *Vec3 {
	return NewVec3(v.e[0]/t, v.e[1]/t, v.e[2]/t)
}

// attributes
func (v *Vec3) Length() float64 {
	return math.Sqrt(v.e[0]*v.e[0] + v.e[1]*v.e[1] + v.e[2]*v.e[2])
}

func (v *Vec3) SquaredLength() float64 {
	return v.e[0]*v.e[0] + v.e[1]*v.e[1] + v.e[2]*v.e[2]
}

func (v *Vec3) Normalize() {
	k := 1.0 / math.Sqrt(v.e[0]*v.e[0]+v.e[1]*v.e[1]+v.e[2]*v.e[2])
	v.e[0] *= k
	v.e[1] *= k
	v.e[2] *= k
}

func (v *Vec3) UnitVector() *Vec3 {
	k := 1.0 / math.Sqrt(v.e[0]*v.e[0]+v.e[1]*v.e[1]+v.e[2]*v.e[2])
	return NewVec3(
		v.e[0]*k,
		v.e[1]*k,
		v.e[2]*k)
}

// dot and cross product
func Dot(v1 *Vec3, v2 *Vec3) float64 {
	return v1.e[0]*v2.e[0] + v1.e[1]*v2.e[1] + v1.e[2]*v2.e[2]
}

func Cross(v1 *Vec3, v2 *Vec3) *Vec3 {
	return NewVec3(
		v1.e[1]*v2.e[2]-v1.e[2]*v2.e[1],
		-(v1.e[0]*v2.e[2] - v1.e[2]*v2.e[0]),
		(v1.e[0]*v2.e[1] - v1.e[1]*v2.e[0]))
}

func RandomCosineDirection() *Vec3 {
	r1 := rand.Float64()
	r2 := rand.Float64()
	z := math.Sqrt(1 - r2)
	phi := 2 * math.Pi * r1
	x := math.Cos(phi) * 2 * math.Sqrt(r2)
	y := math.Sin(phi) * 2 * math.Sqrt(r2)
	return NewVec3(x, y, z)
}
