package geometry

import "math"

type Onb struct {
	Axis [3](*Vec3)
}

func BuildFromW(n *Vec3) *Onb {
	axis := new([3](*Vec3))
	axis[2] = n.UnitVector()
	a := new(Vec3)
	if math.Abs(axis[2].X()) > 0.9 {
		a = NewVec3(0, 1, 0)
	} else {
		a = NewVec3(1, 0, 0)
	}
	axis[1] = Cross(axis[2], a)
	axis[0] = Cross(axis[2], axis[1])
	return &Onb{
		Axis: *axis,
	}
}

func (onb *Onb) U() *Vec3 {
	return onb.Axis[0]
}

func (onb *Onb) V() *Vec3 {
	return onb.Axis[1]
}

func (onb *Onb) W() *Vec3 {
	return onb.Axis[2]
}

func (onb *Onb) Local(a, b, c float64) *Vec3 {
	return (onb.U().TimesScalar(a)).Plus(onb.V().TimesScalar(b)).Plus(onb.W().TimesScalar(c))
}

func (onb *Onb) LocalVector(a *Vec3) *Vec3 {
	return (onb.U().TimesScalar(a.X())).Plus(onb.V().TimesScalar(a.Y())).Plus(onb.W().TimesScalar(a.Z()))
}
