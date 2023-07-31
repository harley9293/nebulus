package util

const F float64 = 1.0 / 0x7fff

type RandGenerator struct {
	seed  int
	count int
}

func NewRandGenerator(seed int) *RandGenerator {
	return &RandGenerator{seed: seed, count: 0}
}

func (r *RandGenerator) RandomInt(max int) int {
	return int(float64(max) * r.RandomFloat())
}

func (r *RandGenerator) RandomFloat() float64 {
	r.count++
	r.seed = r.seed*214013 + 2531011
	val := (float64((r.seed>>16)&0x7fff) - 1) * F
	if val > 0.99999 {
		return 0.99999
	} else {
		return val
	}
}

func (r *RandGenerator) RandomRange(min, max float64) float64 {
	return min + r.RandomFloat()*(max-min)
}
