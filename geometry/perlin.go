package geometry

import (
	"math"
	"math/rand"
)

type Perlin struct {
	RanVec []Vec3
	PermX  []int
	PermY  []int
	PermZ  []int
}

func NewPerlin() *Perlin {
	return &Perlin{
		RanVec: perlinGenerate(),
		PermX:  perlinGeneratePerm(),
		PermY:  perlinGeneratePerm(),
		PermZ:  perlinGeneratePerm(),
	}
}

func (perlin *Perlin) Noise(p *Vec3) float64 {
	u := p.X() - math.Floor(p.X())
	v := p.Y() - math.Floor(p.Y())
	w := p.Z() - math.Floor(p.Z())

	i := int(math.Floor(p.X()))
	j := int(math.Floor(p.Y()))
	k := int(math.Floor(p.Z()))
	c := new([2][2][2]Vec3)
	for di := 0; di < 2; di++ {
		for dj := 0; dj < 2; dj++ {
			for dk := 0; dk < 2; dk++ {
				c[di][dj][dk] = perlin.RanVec[perlin.PermX[(i+di)&255]^perlin.PermY[(j+dj)&255]^perlin.PermZ[(k+dk)&255]]
			}
		}
	}
	return PerlinInterp(c, u, v, w)
}

func (perlin *Perlin) Turb(p *Vec3, depth int) float64 {
	accum := 0.0
	tempP := new(Vec3)
	*tempP = *p
	weight := 1.0
	for i := 0; i < depth; i++ {
		accum += weight * perlin.Noise(tempP)
		weight = weight * 0.5
		*tempP = *tempP.TimesScalar(2)
	}
	return math.Abs(accum)
}

//
func perlinGenerate() []Vec3 {
	p := make([]Vec3, 256)
	for i := 0; i < 256; i++ {
		p[i] = *NewVec3(2*rand.Float64()-1, 2*rand.Float64()-1, 2*rand.Float64()-1).UnitVector()
	}
	return p
}

func permute(p []int, n int) {
	for i := n - 1; i > 0; i-- {
		target := int(rand.Float64() * (float64(i) + 1))
		p[i], p[target] = p[target], p[i]
	}
}

func perlinGeneratePerm() []int {
	p := make([]int, 256)
	for i := 0; i < 256; i++ {
		p[i] = i
	}
	permute(p, 256)
	return p
}

func TrilinearInterp(c *[2][2][2]float64, u, v, w float64) float64 {
	accum := 0.0
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				accum += (float64(i)*u + (1-float64(i))*(1-u)) *
					(float64(j)*v + (1-float64(j))*(1-v)) *
					(float64(k)*w + (1-float64(k))*(1-w)) * c[i][j][k]
			}
		}
	}
	return accum
}

func PerlinInterp(c *[2][2][2]Vec3, u, v, w float64) float64 {
	uu := u * u * (3 - 2*u)
	vv := v * v * (3 - 2*v)
	ww := w * w * (3 - 2*w)
	accum := 0.0
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				weightV := NewVec3(u-float64(i), v-float64(j), w-float64(k))
				newAccum := (float64(i)*uu + (1-float64(i))*(1-uu)) *
					(float64(j)*vv + (1-float64(j))*(1-vv)) *
					(float64(k)*ww + (1-float64(k))*(1-ww)) * 0.5 * (1 - Dot(&c[i][j][k], weightV))
				accum += newAccum
			}
		}
	}

	return accum
}
