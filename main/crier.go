package main

import (
	"log"
	"math"
	"time"
)

type Crier struct {
	HeaderFm
	serveQueue    chan Wish
	dataBase      *Cdb
	inteligents   ObjList
	timeResultion time.Duration
}

func (cr *Crier) Init(dataBase *Cdb) *Crier {
	cr.dataBase = dataBase
	cr.HeaderFm.PostQueue = make(chan Wish, 20)
	cr.serveQueue = make(chan Wish, 50)

	return cr
}

func (cr *Crier) Run() {

	for {
		select {

		case wish, chanOpen := <-cr.serveQueue:
			if !chanOpen {
				log.Printf(" killing Crier")
				return
			}

			if wish["type"] == float64(WISH_PING) {
				wish["receiver"] = 0
				cr.PostQueue <- wish
				continue
			}

			cr.sendSigs(wish)

		// case <-time.After(time.Second / 20):

		case <-time.After(time.Second / 5):
			for iter, intl := cr.dataBase.Inteligents.GetRoot(); iter != nil; iter, intl = cr.dataBase.Inteligents.Next(iter) {

				vel := intl.Brain.Decide()
				log.Printf("Crier >> < %v > decided: %v\n", intl.Id, vel)
				intl.MoveMe(vel[0], vel[1])
			}
		}
	}
}

func (cr *Crier) sendSigs(wish Wish) {

	target, founded := cr.dataBase.MainList.Get(uint(wish["id"]))
	// fmt.Printf("sending %v to <%v>\n", wish["type"], wish["id"])
	if !founded || isNil(target.Brain) {
		return
	}

	sig := new(EnvSig)
	switch wish["type"] {

	case WISH_ALART:
		*sig = EnvSig{
			Type:        SIG_ALART,
			OthersTag:   target.Tag,
			RelPolarPos: [2]float64{wish["radius"], wish["theta"]},
			Reliability: genReliability(wish["radius"], target.Brain.GetSencRadius()),
		}
	case WISH_INTERCOLL:
		*sig = EnvSig{
			Type:        SIG_COLL,
			OthersTag:   target.Tag,
			Reliability: 1,
		}
	}

	log.Printf("Crier >> Sending sig %v to < %v >\n", sig, wish["id"])
	target.Brain.PushEnvSig(sig)

}

func genReliability(dist float64, senc [2]float64) (Reli float64) {
	if dist < senc[0] {
		Reli = 1.0
	} else if dist < senc[1] {
		Reli = 1 - (dist-senc[0])*(0.5/(senc[1]-senc[0]))
	} else {
		math.Max(0, 1-(dist*dist-senc[1])*(0.5/(senc[1])))
	}
	return
}
