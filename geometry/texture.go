package geometry

import "math"

// Interface
type Texture interface {
	Value(u, v float64, p *Vec3) *Vec3
}

// Structs
// Constant
type ConstantTexture struct {
	Color *Vec3
}

func NewConstantTexture(color *Vec3) *ConstantTexture {
	return &ConstantTexture{
		Color: color,
	}
}

func (tex ConstantTexture) Value(u, v float64, p *Vec3) *Vec3 {
	return tex.Color
}

// CheckerTexture
type CheckerTexture struct {
	Odd  Texture
	Even Texture
}

func NewCheckerTexture(t0, t1 Texture) *CheckerTexture {
	return &CheckerTexture{
		Odd:  t1,
		Even: t0,
	}
}

func (tex CheckerTexture) Value(u, v float64, p *Vec3) *Vec3 {
	sines := math.Sin(10*p.X()) * math.Sin(10*p.Y()) * math.Sin(10*p.Z())
	if sines < 0 {
		return tex.Odd.Value(u, v, p)
	}
	return tex.Even.Value(u, v, p)
}

// Noise Texture
type NoiseTexture struct {
	Noise *Perlin
	Scale float64
}

func NewNoiseTexture(scale float64) *NoiseTexture {
	return &NoiseTexture{
		Noise: NewPerlin(),
		Scale: scale,
	}
}

func (tex NoiseTexture) Value(u, v float64, p *Vec3) *Vec3 {
	//return NewVec3(1, 1, 1).TimesScalar(0.5 * (1 + tex.Noise.Turb(p, 7)))
	//return NewVec3(1, 1, 1).TimesScalar(tex.Noise.Turb(p.TimesScalar(tex.Scale), 7))
	return NewVec3(1, 1, 1).TimesScalar(0.5 * (1 + math.Sin(tex.Scale*p.Z()+10*tex.Noise.Turb(p.TimesScalar(tex.Scale), 7))))
}
