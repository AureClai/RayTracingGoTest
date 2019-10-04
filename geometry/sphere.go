package geometry

import (
	"math"
	"math/rand"
)

// Sphere is the type of Sphere
type Sphere struct {
	Center *Vec3
	Radius float64
	Mat    Material
}

func NewSphere(center *Vec3, radius float64, mat Material) *Sphere {
	return &Sphere{center, radius, mat}
}

// Hit test if the sphere is hit
func (sph Sphere) Hit(r *Ray, tMin float64, tMax float64, rec *HitRecord) bool {
	radius := sph.Radius
	var oc = r.Origin().Minus(sph.Center)
	var a = Dot(r.Direction(), r.Direction())
	var b = 2 * Dot(oc, r.Direction())
	var c = Dot(oc, oc) - radius*radius
	var discriminant = b*b - 4*a*c

	if discriminant > 0 {
		temp := (-b - math.Sqrt(b*b-4*a*c)) / (2 * a)
		if temp < tMax && temp > tMin {
			rec.T = temp
			rec.P = r.PointAt(rec.T)
			rec.Normal = rec.P.Minus(sph.Center).ByScalar(radius)
			rec.MatPtr = sph.Mat
			return true
		}
		temp = (-b + math.Sqrt(b*b-4*a*c)) / (2 * a)
		if temp < tMax && temp > tMin {
			rec.T = temp
			rec.P = r.PointAt(rec.T)
			rec.Normal = rec.P.Minus(sph.Center).ByScalar(radius)
			rec.MatPtr = sph.Mat
			return true
		}
	}
	return false
}

func (sph Sphere) BoundingBox(t0, t1 float64, box *Aabb) bool {
	*box = *NewAabb(sph.Center.Minus(NewVec3(sph.Radius, sph.Radius, sph.Radius)), sph.Center.Plus(NewVec3(sph.Radius, sph.Radius, sph.Radius)))
	return true
}

func (sph Sphere) PdfValue(o, v *Vec3) float64 {
	rec := new(HitRecord)
	if sph.Hit(NewRay(o, v), 0.001, math.MaxFloat64, rec) {
		cosThetaMax := math.Sqrt(1 - sph.Radius/(sph.Center.Minus(o)).SquaredLength())
		solidAngle := 2 * math.Pi * (1 - cosThetaMax)
		return 1 / solidAngle
	}
	return 0
}

func (sph Sphere) Random(o *Vec3) *Vec3 {
	direction := sph.Center.Minus(o)
	distanceSquared := direction.SquaredLength()
	uvw := BuildFromW(direction)
	return uvw.LocalVector(randomToSphere(sph.Radius, distanceSquared))
}

func randomInUnitSphere() *Vec3 {
	var p = &Vec3{}
	for {
		p = NewVec3(2*(rand.Float64()-0.5), 2*(rand.Float64()-0.5), 2*(rand.Float64()-0.5))
		if p.SquaredLength() >= 1.0 {
			break
		}
	}
	return p
}

func randomOnUnitSphere() *Vec3 {
	var p = &Vec3{}
	for {
		p = NewVec3(2*(rand.Float64()-0.5), 2*(rand.Float64()-0.5), 2*(rand.Float64()-0.5))
		if p.SquaredLength() >= 1.0 {
			break
		}
	}
	return p.UnitVector()
}

func randomToSphere(radius, distanceSquared float64) *Vec3 {
	r1 := rand.Float64()
	r2 := rand.Float64()
	z := 1 + r2*(math.Sqrt(1-radius*radius/distanceSquared)-1)
	phi := 2 * math.Pi * r1
	x := math.Cos(phi) * math.Sqrt(1-z*z)
	y := math.Sin(phi) * math.Sqrt(1-z*z)
	return NewVec3(x, y, z)
}
