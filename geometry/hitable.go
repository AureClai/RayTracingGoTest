package geometry

import "math/rand"

// HitRecord is the record of the hit
type HitRecord struct {
	T      float64
	U      float64
	V      float64
	P      *Vec3
	Normal *Vec3
	MatPtr Material
}

// Hitable is the interface of all Hitable objects
type Hitable interface {
	Hit(r *Ray, tMin float64, tMax float64, rec *HitRecord) bool
	BoundingBox(t0, t1 float64, box *Aabb) bool
	PdfValue(o, v *Vec3) float64
	Random(o *Vec3) *Vec3
}

// HitableList is a list of Hitable
type HitableList struct {
	list     []Hitable
	listSize int
}

// NewHitableList instanciate a new hitable list
func NewHitableList(list *[]Hitable, listSize int) *HitableList {
	var newList = make([]Hitable, listSize)
	var oldList = *list
	i := 0
	for {
		newList[i] = oldList[i]
		i++
		if i == listSize {
			break
		}
	}
	return &HitableList{newList, listSize}
}

func (hList *HitableList) GetAt(i int) *Hitable {
	return &(hList.list[i])
}

func (hList *HitableList) Slice() (*HitableList, *HitableList) {
	middle := hList.listSize / 2
	firstList := hList.list[:middle]
	firstSize := len(firstList)
	secondList := hList.list[middle:]
	secondSize := len(secondList)
	return NewHitableList(&firstList, firstSize), NewHitableList(&secondList, secondSize)
}

// Hit test if the ray hit any of the hitable in the list
func (hList HitableList) Hit(r *Ray, tMin float64, tMax float64, rec *HitRecord) bool {
	var tempRec = HitRecord{}
	hitAnything := false
	closestSoFar := tMax
	for i := 0; i < hList.listSize; i++ {
		var hitable = hList.GetAt(i)
		if (*hitable).Hit(r, tMin, closestSoFar, &tempRec) {
			hitAnything = true
			closestSoFar = tempRec.T
			*rec = tempRec
		}
	}
	return hitAnything
}

func (hList HitableList) BoundingBox(t0, t1 float64, box *Aabb) bool {
	if hList.listSize < 1 {
		return false
	}
	var tempBox = new(Aabb)
	firstTrue := hList.list[0].BoundingBox(t0, t1, tempBox)
	if !firstTrue {
		return false
	} else {
		*box = *tempBox
	}
	for i := 1; i < hList.listSize; i++ {
		if hList.list[0].BoundingBox(t0, t1, tempBox) {
			*box = *SurroundingBox(box, tempBox)
		} else {
			return false
		}
	}
	return true
}

func (hList HitableList) PdfValue(o, v *Vec3) float64 {
	weight := 1.0 / float64(hList.listSize)
	sum := 0.0
	for i := 0; i < hList.listSize; i++ {
		sum += weight * (*hList.GetAt(i)).PdfValue(o, v)
	}
	return sum
}

func (hList HitableList) Random(o *Vec3) *Vec3 {
	index := int(rand.Float64() * float64(hList.listSize))
	return (*hList.GetAt(index)).Random(o)
}

// Interface sort implementation for liust of Hitable

type byX HitableList
type byY HitableList
type byZ HitableList

// Len
func (hList byX) Len() int {
	return hList.listSize
}
func (hList byY) Len() int {
	return hList.listSize
}
func (hList byZ) Len() int {
	return hList.listSize
}

// Less
func (hList byX) Less(i, j int) bool {
	var boxLeft = new(Aabb)
	var boxRight = new(Aabb)
	if !hList.list[i].BoundingBox(0, 0, boxLeft) || !hList.list[j].BoundingBox(0, 0, boxRight) {
		panic("no bounding box in bvh_node constructor")
	}
	if boxLeft.Min().X() < boxRight.Min().X() {
		return true
	}
	return false
}
func (hList byY) Less(i, j int) bool {
	var boxLeft = new(Aabb)
	var boxRight = new(Aabb)
	if !hList.list[i].BoundingBox(0, 0, boxLeft) || !hList.list[j].BoundingBox(0, 0, boxRight) {
		panic("no bounding box in bvh_node constructor")
	}
	if boxLeft.Min().Y() < boxRight.Min().Y() {
		return true
	}
	return false
}
func (hList byZ) Less(i, j int) bool {
	var boxLeft = new(Aabb)
	var boxRight = new(Aabb)
	if !hList.list[i].BoundingBox(0, 0, boxLeft) || !hList.list[j].BoundingBox(0, 0, boxRight) {
		panic("no bounding box in bvh_node constructor")
	}
	if boxLeft.Min().Z() < boxRight.Min().Z() {
		return true
	}
	return false
}

// Swap
func (hList byX) Swap(i, j int) {
	if i >= hList.listSize || j >= hList.listSize {
		panic("Index out of range")
	}
	hList.list[i], hList.list[j] = hList.list[j], hList.list[i]
}
func (hList byY) Swap(i, j int) {
	if i >= hList.listSize || j >= hList.listSize {
		panic("Index out of range")
	}
	hList.list[i], hList.list[j] = hList.list[j], hList.list[i]
}
func (hList byZ) Swap(i, j int) {
	if i >= hList.listSize || j >= hList.listSize {
		panic("Index out of range")
	}
	hList.list[i], hList.list[j] = hList.list[j], hList.list[i]
}

// Flip normals
type FlipNormals struct {
	Ptr Hitable
}

func NewFlipNormals(ptr Hitable) *FlipNormals {
	return &FlipNormals{
		Ptr: ptr,
	}
}

func (fn FlipNormals) BoundingBox(t0, t1 float64, box *Aabb) bool {
	return fn.Ptr.BoundingBox(t0, t1, box)
}

func (fn FlipNormals) Hit(r *Ray, tMin, tMax float64, rec *HitRecord) bool {
	if fn.Ptr.Hit(r, tMin, tMax, rec) {
		*(rec.Normal) = *(rec.Normal.TimesScalar(-1))
		return true
	}
	return false
}

func (fn FlipNormals) PdfValue(o, v *Vec3) float64 {
	return 0.0
}

func (fn FlipNormals) Random(o *Vec3) *Vec3 {
	return NewVec3(1, 0, 0)
}
