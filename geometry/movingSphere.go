package geometry

import "math"

// Sphere is the type of Sphere
type MovingSphere struct {
	Center0 *Vec3
	Center1 *Vec3
	Radius  float64
	Mat     Material
	Time0   float64
	Time1   float64
}

func NewMovingSphere(cen0, cen1 *Vec3, radius, t0, t1 float64, mat Material) *MovingSphere {
	return &MovingSphere{
		Center0: cen0,
		Center1: cen1,
		Radius:  radius,
		Time0:   t0,
		Time1:   t1,
		Mat:     mat,
	}
}

func (sph *MovingSphere) Center(time float64) *Vec3 {
	return sph.Center0.Plus(sph.Center1.Minus(sph.Center0).TimesScalar((time - sph.Time0) / (sph.Time1 - sph.Time0)))
}

// Hit test if the sphere is hit
func (sph MovingSphere) Hit(r *Ray, tMin float64, tMax float64, rec *HitRecord) bool {
	var oc = r.Origin().Minus(sph.Center(r.Time()))
	var a = Dot(r.Direction(), r.Direction())
	var b = 2 * Dot(oc, r.Direction())
	var c = Dot(oc, oc) - sph.Radius*sph.Radius
	var discriminant = b*b - 4*a*c

	if discriminant > 0 {
		temp := (-b - math.Sqrt(b*b-4*a*c)) / (2 * a)
		if temp < tMax && temp > tMin {
			rec.T = temp
			rec.P = r.PointAt(rec.T)
			rec.Normal = rec.P.Minus(sph.Center(r.Time())).ByScalar(sph.Radius)
			rec.MatPtr = sph.Mat
			return true
		}
		temp = (-b + math.Sqrt(b*b-4*a*c)) / (2 * a)
		if temp < tMax && temp > tMin {
			rec.T = temp
			rec.P = r.PointAt(rec.T)
			rec.Normal = rec.P.Minus(sph.Center(r.Time())).ByScalar(sph.Radius)
			rec.MatPtr = sph.Mat
			return true
		}
	}
	return false
}

func (sph MovingSphere) BoundingBox(t0, t1 float64, box *Aabb) bool {
	var box0 = NewAabb(sph.Center(t0).Minus(NewVec3(sph.Radius, sph.Radius, sph.Radius)), sph.Center(t0).Plus(NewVec3(sph.Radius, sph.Radius, sph.Radius)))
	var box1 = NewAabb(sph.Center(t1).Minus(NewVec3(sph.Radius, sph.Radius, sph.Radius)), sph.Center(t1).Plus(NewVec3(sph.Radius, sph.Radius, sph.Radius)))
	*box = *SurroundingBox(box0, box1)
	return true
}
