package main

import "math"

const PI = 3.14

type GeoShape struct {
	Typee      string
	Properties []float64
}

func distance(p1, p2 [2]float64, isPolar bool) float64 {
	if isPolar {
		return math.Abs(p1[0] - p2[0])
	}
	return math.Sqrt(math.Pow(p1[0]-p2[0], 2) + math.Pow(p1[1]-p2[1], 2))
}

func cartesianToPolar(cart [2]float64) (pol [2]float64) {
	if cart == [2]float64{0, 0} {
		return cart
	}

	pol[0] = distance(cart, [2]float64{0, 0}, false)
	pol[1] = math.Atan(cart[1] / cart[0])
	return
}

func polarToCartasian(pol [2]float64) (cart [2]float64) {
	cart[0] = pol[0] * math.Cos(pol[1])
	cart[1] = pol[0] * math.Sin(pol[1])
	return
}

func vecDo(vec1, vec2 []float64, operation func(float64, float64) float64) (res []float64) {
	if len(vec1) != len(vec2) {
		panic("vecs with different length")
	}

	for i := 0; i < len(vec1); i++ {
		res[i] = operation(vec1[i], vec2[i])
	}
	return
}

func vecMinus(vec1, vec2 []float64) []float64 {
	return vecDo(vec1, vec2, func(e1, e2 float64) float64 { return e1 - e2 })
}
