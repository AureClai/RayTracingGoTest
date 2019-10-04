package geometry

import (
	"sort"
)

type BVHNode struct {
	Left  Hitable
	Right Hitable
	Box   *Aabb
}

func NewBVHNode(l *HitableList, time0, time1 float64) *BVHNode {
	axis := int(3 * drand48())
	if axis == 0 {
		sort.Sort(byX(*l))
	} else if axis == 1 {
		sort.Sort(byY(*l))
	} else {
		sort.Sort(byZ(*l))
	}
	var left = new(Hitable)
	var right = new(Hitable)
	if l.listSize == 1 {
		*left = l.list[0]
		*right = l.list[0]
	} else if l.listSize == 2 {
		*left = l.list[0]
		*right = l.list[1]
	} else {
		newHList0, newHList1 := l.Slice()
		*left = NewBVHNode(newHList0, time0, time1)
		*right = NewBVHNode(newHList1, time0, time1)
	}
	var boxLeft = new(Aabb)
	var boxRight = new(Aabb)
	if !(*left).BoundingBox(time0, time1, boxLeft) || !(*right).BoundingBox(time0, time1, boxRight) {
		panic("no bounding box in bvh_node constructor")
	}
	box := SurroundingBox(boxLeft, boxRight)
	return &BVHNode{
		Left:  *left,
		Right: *right,
		Box:   box,
	}
}

func (bvhn BVHNode) Hit(r *Ray, tMin float64, tMax float64, rec *HitRecord) bool {
	if bvhn.Box.Hit(r, tMin, tMax) {
		var leftRec = new(HitRecord)
		var rightRec = new(HitRecord)
		hitLeft := bvhn.Left.Hit(r, tMin, tMax, leftRec)
		hitRight := bvhn.Right.Hit(r, tMin, tMax, rightRec)
		if hitLeft && hitRight {
			if leftRec.T < rightRec.T {
				*rec = *leftRec
			} else {
				*rec = *rightRec
			}
			return true
		} else if hitLeft {
			*rec = *leftRec
			return true
		} else if hitRight {
			*rec = *rightRec
			return true
		} else {
			return false
		}
	}
	return false
}

func (bvhn BVHNode) BoundingBox(t0, t1 float64, b *Aabb) bool {
	*b = *bvhn.Box
	return true
}

func (bvhn BVHNode) PdfValue(o, v *Vec3) float64 {
	return 0.0
}

func (bvhn BVHNode) Random(o *Vec3) *Vec3 {
	return NewVec3(1, 0, 0)
}
