package geometry

import (
	"math"
	"math/rand"
)

/*
* Utilitaries
 */
func reflect(v *Vec3, n *Vec3) *Vec3 {
	return v.Minus(n.TimesScalar(2 * Dot(v, n)))
}

func refract(v *Vec3, n *Vec3, niOverNt float64, refracted *Vec3) bool {
	var uv = v.UnitVector()
	dt := Dot(uv, n)
	discriminant := 1.0 - niOverNt*niOverNt*(1-dt*dt)
	if discriminant > 0 {
		*refracted = *uv.Minus(n.TimesScalar(dt)).TimesScalar(niOverNt).Minus(n.TimesScalar(math.Sqrt(discriminant)))
		return true
	}
	return false

}

func schlick(cosine, refIdx float64) float64 {
	r0 := (1 - refIdx) / (1 + refIdx)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}

/*
* 	Interface for materials
 */

// Material interface
// struct of scattered record
type ScatterRecord struct {
	SpecularRay *Ray
	IsSpecular  bool
	Attenuation *Vec3
	PdfPtr      Pdf
}

type Material interface {
	Scatter(rIn *Ray, rec *HitRecord, srec *ScatterRecord) bool
	Emitted(rIn *Ray, rec *HitRecord, u, v float64, p *Vec3) *Vec3
	ScatteringPdf(rIn *Ray, rec *HitRecord, scattered *Ray) float64
}

type noMaterial struct{}

func NewNoMaterial() *noMaterial {
	return &noMaterial{}
}

func (noMat *noMaterial) Scatter(rIn *Ray, hrec *HitRecord, srec *ScatterRecord) bool {
	return false
}

func (noMat *noMaterial) Emitted(rIn *Ray, rec *HitRecord, u, v float64, p *Vec3) *Vec3 {
	return NewVec3(0, 0, 0)
}

func (noMat *noMaterial) ScatteringPdf(rIn *Ray, rec *HitRecord, scattered *Ray) float64 {
	return 0.0
}

/*
* MATERIALS AND FUNCTION
 */

// Lambertian is a material
type Lambertian struct {
	Albedo Texture
}

func (lamb Lambertian) Scatter(rIn *Ray, hrec *HitRecord, srec *ScatterRecord) bool {
	srec.IsSpecular = false
	srec.Attenuation = lamb.Albedo.Value(hrec.U, hrec.V, hrec.P)
	srec.PdfPtr = NewCosinePdf(hrec.Normal)
	return true
}

func (lamb Lambertian) Emitted(rIn *Ray, rec *HitRecord, u, v float64, p *Vec3) *Vec3 {
	return NewVec3(0, 0, 0)
}

func (lamb Lambertian) ScatteringPdf(rIn *Ray, rec *HitRecord, scattered *Ray) float64 {
	cosine := Dot(rec.Normal, scattered.Direction().UnitVector())
	if cosine < 0 {
		cosine = 0
	}
	return cosine / math.Pi
}

// Metal material
type Metal struct {
	Albedo *Vec3
	Fuzz   float64
}

func (met Metal) Scatter(rIn *Ray, hrec *HitRecord, srec *ScatterRecord) bool {
	var reflected = reflect(rIn.Direction().UnitVector(), hrec.Normal)
	srec.SpecularRay = NewRay(hrec.P, reflected.Plus(randomInUnitSphere().TimesScalar(met.Fuzz)))
	srec.Attenuation = met.Albedo
	srec.IsSpecular = true
	srec.PdfPtr = NewNoPdf()
	return true
}

func (met Metal) Emitted(rIn *Ray, rec *HitRecord, u, v float64, p *Vec3) *Vec3 {
	return NewVec3(0, 0, 0)
}

func (met Metal) ScatteringPdf(rIn *Ray, rec *HitRecord, scattered *Ray) float64 {
	return 0.0
}

// Dielectric material
type Dielectric struct {
	RefIdx float64
}

func (die Dielectric) Scatter(rIn *Ray, hrec *HitRecord, srec *ScatterRecord) bool {
	srec.IsSpecular = true
	srec.PdfPtr = NewNoPdf()
	srec.Attenuation = NewVec3(1.0, 1.0, 1.0)

	outwardNormal := new(Vec3)
	reflected := reflect(rIn.Direction(), hrec.Normal)
	refracted := new(Vec3)

	niOverNt := 0.0
	reflectProb := 0.0
	cosine := 0.0

	if Dot(rIn.Direction(), hrec.Normal) > 0 {
		outwardNormal = hrec.Normal.Opposite()
		niOverNt = die.RefIdx
		cosine = die.RefIdx * Dot(rIn.Direction(), hrec.Normal) / rIn.Direction().Length()
	} else {
		outwardNormal = hrec.Normal
		niOverNt = 1.0 / die.RefIdx
		cosine = -Dot(rIn.Direction(), hrec.Normal) / rIn.Direction().Length()
	}
	if refract(rIn.Direction(), outwardNormal, niOverNt, refracted) {
		reflectProb = schlick(cosine, die.RefIdx)
	} else {
		reflectProb = 1.0
	}
	if rand.Float64() < reflectProb {
		srec.SpecularRay = NewRay(hrec.P, reflected)
	} else {
		srec.SpecularRay = NewRay(hrec.P, refracted)
	}
	return true
}

func (die Dielectric) Emitted(rIn *Ray, rec *HitRecord, u, v float64, p *Vec3) *Vec3 {
	return NewVec3(0, 0, 0)
}

func (die Dielectric) ScatteringPdf(rIn *Ray, rec *HitRecord, scattered *Ray) float64 {
	return 0.0
}

// Diffuse Light
type DiffuseLight struct {
	Emit Texture
}

func NewDiffuseLight(emit Texture) *DiffuseLight {
	return &DiffuseLight{
		Emit: emit,
	}
}

func (dl DiffuseLight) Scatter(rIn *Ray, hrec *HitRecord, srec *ScatterRecord) bool {
	return false
}

func (dl DiffuseLight) ScatteringPdf(rIn *Ray, rec *HitRecord, scattered *Ray) float64 {
	return 0.0
}

func (dl DiffuseLight) Emitted(rIn *Ray, rec *HitRecord, u, v float64, p *Vec3) *Vec3 {
	if Dot(rec.Normal, rIn.Direction()) < 0.0 {
		return dl.Emit.Value(u, v, p)
	}
	return NewVec3(0, 0, 0)
}

//

type Isotropic struct {
	Albedo Texture
}

func NewIsotropic(a Texture) *Isotropic {
	return &Isotropic{
		Albedo: a,
	}
}

func (iso Isotropic) Scatter(rIn *Ray, rec *HitRecord, attenuation *Vec3, scattered *Ray) bool {
	*scattered = *NewRay(rec.P, randomInUnitSphere())
	*attenuation = *iso.Albedo.Value(rec.U, rec.V, rec.P)
	return true
}

func (iso Isotropic) Emitted(u, v float64, p *Vec3) *Vec3 {
	return NewVec3(0, 0, 0)
}
