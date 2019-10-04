package geometry

import (
	"math"
	"math/rand"
)

type Pdf interface {
	Value(direction *Vec3) float64
	Generate() *Vec3
}

//
type noPdf struct{}

func NewNoPdf() *noPdf {
	return &noPdf{}
}

func (cpdf noPdf) Value(direction *Vec3) float64 {
	return 0.0
}

func (cpdf noPdf) Generate() *Vec3 {
	return NewVec3(0, 0, 0)
}

//

type CosinePdf struct {
	Uvw *Onb
}

func NewCosinePdf(w *Vec3) *CosinePdf {
	return &CosinePdf{
		Uvw: BuildFromW(w),
	}
}

func (cpdf CosinePdf) Value(direction *Vec3) float64 {
	cosine := Dot(direction.UnitVector(), cpdf.Uvw.W())
	if cosine > 0 {
		return cosine / math.Pi
	}
	return 0.0
}

func (cpdf CosinePdf) Generate() *Vec3 {
	return cpdf.Uvw.LocalVector(RandomCosineDirection())
}

//
type HitablePdf struct {
	O   *Vec3
	Ptr Hitable
}

func NewHitablePdf(p Hitable, origin *Vec3) *HitablePdf {
	return &HitablePdf{
		O:   origin,
		Ptr: p,
	}
}

func (hpdf HitablePdf) Value(direction *Vec3) float64 {
	return hpdf.Ptr.PdfValue(hpdf.O, direction)
}

func (hpdf HitablePdf) Generate() *Vec3 {
	return hpdf.Ptr.Random(hpdf.O)
}

//
type MixturePdf struct {
	P [2]Pdf
}

func NewMixturePdf(p0, p1 Pdf) *MixturePdf {
	return &MixturePdf{
		P: [2]Pdf{p0, p1},
	}
}

func (mpdf MixturePdf) Value(direction *Vec3) float64 {
	return 0.5*mpdf.P[0].Value(direction) + 0.5*mpdf.P[1].Value(direction)
}

func (mpdf MixturePdf) Generate() *Vec3 {
	if rand.Float64() < 0.5 {
		return mpdf.P[0].Generate()
	}
	return mpdf.P[1].Generate()
}
