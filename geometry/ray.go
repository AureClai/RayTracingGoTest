package geometry

// Ray is the ray of light in the scene
type Ray struct {
	a     *Vec3
	b     *Vec3
	_time float64
}

// NewRay instantiate a new Ray and return the pointer
func NewRay(_a, _b *Vec3) *Ray {
	return &Ray{a: _a, b: _b, _time: 0.0}
}

func NewRayWithTime(_a, _b *Vec3, ti float64) *Ray {
	return &Ray{a: _a, b: _b, _time: ti}
}

func NewEmptyRay() *Ray {
	return NewRay(NewVec3(0, 0, 0), NewVec3(0, 0, 0))
}

// Origin get the origin point of the Ray
func (r *Ray) Origin() *Vec3 {
	return r.a
}

// Direction get the Vector of direction of the Ray
func (r *Ray) Direction() *Vec3 {
	return r.b
}

func (r *Ray) Time() float64 {
	return r._time
}

// PointAt get the point from A to B at time t
func (r *Ray) PointAt(t float64) *Vec3 {
	return r.a.Plus(r.b.TimesScalar(t))
}
