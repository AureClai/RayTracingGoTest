package geometry

import "math"

type Translate struct {
	Ptr    Hitable
	Offset *Vec3
}

func NewTranslate(p Hitable, displacement *Vec3) *Translate {
	return &Translate{
		Ptr:    p,
		Offset: displacement,
	}
}

func (tr Translate) BoundingBox(t0, t1 float64, box *Aabb) bool {
	if tr.Ptr.BoundingBox(t0, t1, box) {
		*box = *NewAabb(box.Min().Plus(tr.Offset), box.Max().Plus(tr.Offset))
		return true
	}
	return false
}

func (tr Translate) Hit(r *Ray, tMin, tMax float64, rec *HitRecord) bool {
	movedR := NewRayWithTime(r.Origin().Minus(tr.Offset), r.Direction(), r.Time())
	if tr.Ptr.Hit(movedR, tMin, tMax, rec) {
		*rec.P = *(rec.P.Plus(tr.Offset))
		return true
	}
	return false
}

func (tr Translate) PdfValue(o, v *Vec3) float64 {
	return 0.0
}

func (tr Translate) Random(o *Vec3) *Vec3 {
	return NewVec3(1, 0, 0)
}

//

type RotateY struct {
	Ptr      Hitable
	SinTheta float64
	CosTheta float64
	HasBox   bool
	Bbox     *Aabb
}

func NewRotateY(p Hitable, angle float64) *RotateY {
	radians := (math.Pi / 180.0) * angle
	sinTheta := math.Sin(radians)
	cosTheta := math.Cos(radians)
	bbox := new(Aabb)
	hasBox := p.BoundingBox(0, 1, bbox)
	min := NewVec3(math.MaxFloat64, math.MaxFloat64, math.MaxFloat64)
	max := NewVec3(-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64)
	for i := 0; i < 2; i++ {
		x := float64(i)*bbox.Max().X() + (float64(i)-1)*bbox.Min().X()
		for j := 0; j < 2; j++ {
			y := float64(j)*bbox.Max().Y() + (float64(j)-1)*bbox.Min().Y()
			for k := 0; k < 2; k++ {
				z := float64(k)*bbox.Max().Z() + (float64(k)-1)*bbox.Min().Z()
				newX := cosTheta*x + sinTheta*z
				newZ := -sinTheta*x + cosTheta*z
				tester := NewVec3(newX, y, newZ)
				for c := 0; c < 3; c++ {
					if tester.At(c) > max.At(c) {
						max.SetAt(c, tester.At(c))
					}
					if tester.At(c) < min.At(c) {
						min.SetAt(c, tester.At(c))
					}
				}
			}
		}
	}
	return &RotateY{
		Ptr:      p,
		SinTheta: sinTheta,
		CosTheta: cosTheta,
		HasBox:   hasBox,
		Bbox:     bbox,
	}
}

func (ry RotateY) BoundingBox(t0, t1 float64, box *Aabb) bool {
	*box = *ry.Bbox
	return ry.HasBox
}

func (ry RotateY) Hit(r *Ray, tMin, tMax float64, rec *HitRecord) bool {
	origin := new(Vec3)
	*origin = *r.Origin()
	direction := new(Vec3)
	*direction = *r.Direction()
	origin.SetAt(0, ry.CosTheta*r.Origin().At(0)-ry.SinTheta*r.Origin().At(2))
	origin.SetAt(2, ry.SinTheta*r.Origin().At(0)+ry.CosTheta*r.Origin().At(2))
	direction.SetAt(0, ry.CosTheta*r.Direction().At(0)-ry.SinTheta*r.Direction().At(2))
	direction.SetAt(2, ry.SinTheta*r.Direction().At(0)+ry.CosTheta*r.Direction().At(2))
	rotatedR := NewRayWithTime(origin, direction, r.Time())
	if ry.Ptr.Hit(rotatedR, tMin, tMax, rec) {
		p := new(Vec3)
		*p = *rec.P
		normal := new(Vec3)
		*normal = *rec.Normal
		p.SetAt(0, ry.CosTheta*rec.P.At(0)+ry.SinTheta*rec.P.At(2))
		p.SetAt(2, -ry.SinTheta*rec.P.At(0)+ry.CosTheta*rec.P.At(2))
		normal.SetAt(0, ry.CosTheta*rec.Normal.At(0)+ry.SinTheta*rec.Normal.At(2))
		normal.SetAt(2, -ry.SinTheta*rec.Normal.At(0)+ry.CosTheta*rec.Normal.At(2))
		*rec.P = *p
		*rec.Normal = *normal
		return true
	}
	return false
}

func (ry RotateY) PdfValue(o, v *Vec3) float64 {
	return 0.0
}

func (ry RotateY) Random(o *Vec3) *Vec3 {
	return NewVec3(1, 0, 0)
}
