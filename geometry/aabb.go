package geometry

/*
func ffmin(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func ffmax(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
*/

type Aabb struct {
	_min *Vec3
	_max *Vec3
}

func NewAabb(a, b *Vec3) *Aabb {
	return &Aabb{
		_min: a,
		_max: b,
	}
}

func (aabb *Aabb) Min() *Vec3 {
	return aabb._min
}

func (aabb *Aabb) Max() *Vec3 {
	return aabb._max
}

func SurroundingBox(box0, box1 *Aabb) *Aabb {
	var small = NewVec3(
		Fmin(box0.Min().X(), box1.Min().X()),
		Fmin(box0.Min().Y(), box1.Min().Y()),
		Fmin(box0.Min().Z(), box1.Min().Z()))
	var big = NewVec3(
		Fmin(box0.Max().X(), box1.Max().X()),
		Fmin(box0.Max().Y(), box1.Max().Y()),
		Fmin(box0.Max().Z(), box1.Max().Z()))
	return NewAabb(small, big)
}

func (aabb *Aabb) Hit(r *Ray, tMin, tMax float64) bool {
	for a := 0; a < 3; a++ {
		invD := 1.0 / r.Direction().At(a)
		t0 := (aabb._min.At(a) - r.Origin().At(a)) * invD
		t1 := (aabb._max.At(a) - r.Origin().At(a)) * invD
		if invD < 0.0 {
			t0, t1 = t1, t0
		}
		if t0 > tMin {
			tMin = t0
		}
		if t1 < tMax {
			tMax = t1
		}
		if tMax <= tMin {
			return false
		}
	}
	return true
}
