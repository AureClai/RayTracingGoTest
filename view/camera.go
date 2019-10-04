package view

import (
	"math"
	"math/rand"

	g "../geometry"
)

func randomInUnitDisk() *g.Vec3 {
	var p = new(g.Vec3)

	for {
		p = g.NewVec3(2*(rand.Float64()-0.5), 2*(rand.Float64()-0.5), 2*(rand.Float64()-0.5))
		if g.Dot(p, p) >= 1 {
			break
		}
	}
	return p
}

// Camera is a struct that modelize the camera of the scene
type Camera struct {
	Origin          *g.Vec3
	Horizontal      *g.Vec3
	Vertical        *g.Vec3
	LowerLeftCorner *g.Vec3
	U               *g.Vec3
	V               *g.Vec3
	W               *g.Vec3
	LensRadius      float64
	Time0           float64
	Time1           float64
}

func NewCamera(lookFrom, lookAt, vup *g.Vec3, vfov, aspect, aperture, focusDist, t0, t1 float64) *Camera {
	lensRadius := aperture / 2
	theta := vfov * math.Pi / 180
	halfHeight := math.Tan(theta / 2)
	halfWidth := aspect * halfHeight
	var origin = lookFrom
	var w = lookAt.Minus(lookFrom).UnitVector()
	var u = g.Cross(vup, w).UnitVector()
	var v = g.Cross(w, u)
	var lowerLeftCorner = g.NewVec3(-halfWidth, -halfHeight, -1.0)
	lowerLeftCorner = origin.Minus(u.TimesScalar(halfWidth * focusDist)).Minus(v.TimesScalar(halfHeight * focusDist).Minus(w.TimesScalar(focusDist)))
	var horizontal = u.TimesScalar(2 * halfWidth * focusDist)
	var vertical = v.TimesScalar(2 * halfHeight * focusDist)
	return &Camera{Origin: origin,
		Horizontal:      horizontal,
		Vertical:        vertical,
		LowerLeftCorner: lowerLeftCorner,
		U:               u,
		V:               v,
		W:               w,
		LensRadius:      lensRadius,
		Time0:           t0,
		Time1:           t1,
	}
}

func (cam *Camera) GetRay(s, t float64) *g.Ray {
	var rd = randomInUnitDisk().TimesScalar(cam.LensRadius)
	var offset = cam.U.TimesScalar(rd.X()).Plus(cam.V.TimesScalar(rd.Y()))
	var time = cam.Time0 + rand.Float64()*(cam.Time1-cam.Time0)
	return g.NewRayWithTime(cam.Origin.Plus(offset), cam.LowerLeftCorner.Plus(cam.Horizontal.TimesScalar(s)).Plus(cam.Vertical.TimesScalar(t)).Minus(cam.Origin).Minus(offset), time)
}
