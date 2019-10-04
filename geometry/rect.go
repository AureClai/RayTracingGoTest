package geometry

import (
	"math"
	"math/rand"
)

//TODO: Complete Hitable funcs to YZ and XY

type XYRect struct {
	Mat Material
	X0  float64
	X1  float64
	Y0  float64
	Y1  float64
	K   float64
}

func NewXYRect(_x0, _x1, _y0, _y1, _k float64, _mat Material) *XYRect {
	return &XYRect{
		X0:  _x0,
		X1:  _x1,
		Y0:  _y0,
		Y1:  _y1,
		K:   _k,
		Mat: _mat,
	}
}

func (rect XYRect) BoundingBox(t0, t1 float64, box *Aabb) bool {
	*box = *NewAabb(NewVec3(rect.X0, rect.Y0, rect.K-0.0001), NewVec3(rect.X1, rect.Y1, rect.K+0.0001))
	return true
}

func (rect XYRect) Hit(r *Ray, t0, t1 float64, rec *HitRecord) bool {
	t := (rect.K - r.Origin().Z()) / r.Direction().Z()
	if t < t0 || t > t1 {
		return false
	}
	x := r.Origin().X() + t*r.Direction().X()
	y := r.Origin().Y() + t*r.Direction().Y()
	if x < rect.X0 || x > rect.X1 || y < rect.Y0 || y > rect.Y1 {
		return false
	}
	rec.U = (x - rect.X0) / (rect.X1 - rect.X0)
	rec.V = (y - rect.Y0) / (rect.Y1 - rect.Y0)
	rec.T = t
	rec.MatPtr = rect.Mat
	rec.P = r.PointAt(t)
	// x^y = z
	rec.Normal = Cross(NewVec3(rect.X1-rect.X0, 0, 0), NewVec3(0, rect.Y1-rect.Y0, 0)).UnitVector()
	return true
}

func (rect XYRect) PdfValue(o, v *Vec3) float64 {
	return 0.0
}

func (rect XYRect) Random(o *Vec3) *Vec3 {
	return NewVec3(1, 0, 0)
}

//
type XZRect struct {
	Mat Material
	X0  float64
	X1  float64
	Z0  float64
	Z1  float64
	K   float64
}

func NewXZRect(_x0, _x1, _z0, _z1, _k float64, _mat Material) *XZRect {
	return &XZRect{
		X0:  _x0,
		X1:  _x1,
		Z0:  _z0,
		Z1:  _z1,
		K:   _k,
		Mat: _mat,
	}
}

func (rect XZRect) BoundingBox(t0, t1 float64, box *Aabb) bool {
	*box = *NewAabb(NewVec3(rect.X0, rect.K-0.0001, rect.Z0), NewVec3(rect.X1, rect.K+0.0001, rect.Z1))
	return true
}

func (rect XZRect) Hit(r *Ray, t0, t1 float64, rec *HitRecord) bool {
	t := (rect.K - r.Origin().Y()) / r.Direction().Y()
	if t < t0 || t > t1 {
		return false
	}
	x := r.Origin().X() + t*r.Direction().X()
	z := r.Origin().Z() + t*r.Direction().Z()
	if x < rect.X0 || x > rect.X1 || z < rect.Z0 || z > rect.Z1 {
		return false
	}
	rec.U = (x - rect.X0) / (rect.X1 - rect.X0)
	rec.V = (z - rect.Z0) / (rect.Z1 - rect.Z0)
	rec.T = t
	rec.MatPtr = rect.Mat
	rec.P = r.PointAt(t)
	// z^x = y
	rec.Normal = Cross(NewVec3(0, 0, rect.Z1-rect.Z0), NewVec3(rect.X1-rect.X0, 0, 0)).UnitVector()
	return true
}

func (rect XZRect) PdfValue(o, v *Vec3) float64 {
	rec := new(HitRecord)
	if rect.Hit(NewRay(o, v), 0.001, math.MaxFloat64, rec) {
		area := (rect.X1 - rect.X0) * (rect.Z1 - rect.Z0)
		distanceSquared := rec.T * rec.T * v.SquaredLength()
		cosine := math.Abs(Dot(v, rec.Normal) / v.Length())
		return distanceSquared / (cosine * area)
	}
	return 0.0
}

func (rect XZRect) Random(o *Vec3) *Vec3 {
	randomPoint := NewVec3(rect.X0+rand.Float64()*(rect.X1-rect.X0), rect.K, rect.Z0+rand.Float64()*(rect.Z1-rect.Z0))
	return randomPoint.Minus(o)
}

//
type YZRect struct {
	Mat Material
	Y0  float64
	Y1  float64
	Z0  float64
	Z1  float64
	K   float64
}

func NewYZRect(_y0, _y1, _z0, _z1, _k float64, _mat Material) *YZRect {
	return &YZRect{
		Y0:  _y0,
		Y1:  _y1,
		Z0:  _z0,
		Z1:  _z1,
		K:   _k,
		Mat: _mat,
	}
}

func (rect YZRect) BoundingBox(t0, t1 float64, box *Aabb) bool {
	*box = *NewAabb(NewVec3(rect.K-0.0001, rect.Y0, rect.Z0), NewVec3(rect.K+0.0001, rect.Y1, rect.Z1))
	return true
}

func (rect YZRect) Hit(r *Ray, t0, t1 float64, rec *HitRecord) bool {
	t := (rect.K - r.Origin().X()) / r.Direction().X()
	if t < t0 || t > t1 {
		return false
	}
	y := r.Origin().Y() + t*r.Direction().Y()
	z := r.Origin().Z() + t*r.Direction().Z()
	if y < rect.Y0 || y > rect.Y1 || z < rect.Z0 || z > rect.Z1 {
		return false
	}
	rec.U = (y - rect.Y0) / (rect.Y1 - rect.Y0)
	rec.V = (z - rect.Z0) / (rect.Z1 - rect.Z0)
	rec.T = t
	rec.MatPtr = rect.Mat
	rec.P = r.PointAt(t)
	// y^z = x
	rec.Normal = Cross(NewVec3(0, rect.Y1-rect.Y0, 0), NewVec3(0, 0, rect.Z1-rect.Z0)).UnitVector()
	return true
}

func (rect YZRect) PdfValue(o, v *Vec3) float64 {
	return 0.0
}

func (rect YZRect) Random(o *Vec3) *Vec3 {
	return NewVec3(1, 0, 0)
}
