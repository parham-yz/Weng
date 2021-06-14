package main

import (
	"log"
	"math/rand"
)

type gongishkBrain struct {
	dVel [2]float64
}

func (g *gongishkBrain) Decide() [2]float64 {
	log.Printf("loog >> %v\n")
	if g.dVel[0] == 0 {
		// log.Println([2]float64{1, rand.Float64() * PI / 2})
		return [2]float64{1, rand.Float64() * PI / 2}
	}
	return g.dVel
}

func (g *gongishkBrain) PushEnvSig(sig *EnvSig) {

	g.dVel = [2]float64{sig.RelPolarPos[0] * sig.Reliability, sig.RelPolarPos[1] + PI}
}

func (g *gongishkBrain) GetSencRadius() [2]float64 {
	return [2]float64{3, 4}
}
