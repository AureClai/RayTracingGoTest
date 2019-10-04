package geometry

type Box struct {
	PMin    *Vec3
	PMax    *Vec3
	ListPtr *HitableList
}

func NewBox(p0, p1 *Vec3, ptr Material) *Box {
	list := make([]Hitable, 6)
	list[0] = NewXYRect(p0.X(), p1.X(), p0.Y(), p1.Y(), p1.Z(), ptr)
	list[1] = NewFlipNormals(NewXYRect(p0.X(), p1.X(), p0.Y(), p1.Y(), p0.Z(), ptr))
	list[2] = NewXZRect(p0.X(), p1.X(), p0.Z(), p1.Z(), p1.Y(), ptr)
	list[3] = NewFlipNormals(NewXZRect(p0.X(), p1.X(), p0.Z(), p1.Z(), p0.Y(), ptr))
	list[4] = NewYZRect(p0.Y(), p1.Y(), p0.Z(), p1.Z(), p1.X(), ptr)
	list[5] = NewFlipNormals(NewYZRect(p0.Y(), p1.Y(), p0.Z(), p1.Z(), p0.X(), ptr))
	return &Box{
		PMin:    p0,
		PMax:    p1,
		ListPtr: NewHitableList(&list, 6),
	}
}

func (bx Box) BoundingBox(t0, t1 float64, box *Aabb) bool {
	*box = *NewAabb(bx.PMin, bx.PMax)
	return true
}

func (bx Box) Hit(r *Ray, t0, t1 float64, rec *HitRecord) bool {
	return bx.ListPtr.Hit(r, t0, t1, rec)
}

func (bx Box) PdfValue(o, v *Vec3) float64 {
	return 0.0
}

func (bx Box) Random(o *Vec3) *Vec3 {
	return NewVec3(1, 0, 0)
}
