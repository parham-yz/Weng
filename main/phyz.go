package main

import (
	"log"
	"math"
	"time"
)

const (
	BOUNDARY_STATIC = iota
	BOUNDARY_CONTINOUOUS
)

type Phyz struct {
	HeaderFm
	serveQueue       chan Wish
	dataBase         *Cdb
	objBank          *ObjList
	woldProps        map[string]float64
	moversCache      *ObjList
	interestingColls *ObjList
	alive            bool
	timeResultion    time.Duration
}

// Init !!!
func (ph *Phyz) Init(cdb *Cdb, woldProps map[string]float64) *Phyz {
	ph.PostQueue = make(chan Wish, 100)
	ph.serveQueue = make(chan Wish, 200)

	ph.dataBase = cdb
	ph.objBank = &cdb.MainList
	ph.moversCache = &cdb.Movers
	ph.interestingColls = &cdb.Inteligents

	ph.woldProps = woldProps
	ph.Id = 1
	ph.timeResultion = time.Second / 20
	return ph
}

// spwans on object
func (ph *Phyz) spwan(target *Object) bool {

	if target.Phfm == nil {
		return false
	}

	target.Phfm.Present = true
	iss := ph.CheckCollision(target)
	if iss {
		target.Phfm.Present = false
		return false
	}
	log.Printf("Phyz >> <%v> was spawned\n", target.Id)
	return true
}

// Run !!!
func (ph *Phyz) Run() {

	ph.alive = true
	go ph.spinn()

	//serving wishes
	for {

		select {
		case wish, chanOpen := <-ph.serveQueue:
			if !chanOpen {
				log.Println("killing Phyz")
				ph.alive = false
				return
			}
			ph.serveWish(wish)
		case <-time.After(time.Second / 5):
			ph.checkDistances()
		}

	}
}

func (ph *Phyz) serveWish(wish Wish) {

	target, founded := ph.objBank.Get(uint(wish["id"]))
	if !founded {

		log.Printf("phyz >> serving wish <%v>: can not find obj wiht id: %v in obj_bank\n", wish, uint(wish["id"]))
		return
	}

	if target.Phfm == nil {
		log.Printf("phyz >> serving wish <%v>: it is not a phyzical being!\n", wish)
		return
	}

	// switching
	switch wish["type"] {

	case WISH_MOVE:

		if target.Phfm.Present {
			target.Phfm.Velocity[0], target.Phfm.Velocity[1] = wish["radius"], wish["theta"]

			if target.Phfm.Velocity[0] == 0 && target.Phfm.Velocity[1] == 0 {
				ph.moversCache.Del(target.Id)
				// log.Printf("Phyz >> <%v> is not a Mover any more\n", target.Id)
			} else {
				ph.moversCache.Add(target)
				// log.Printf("Phyz >> <%v> is a Mover: (%v)\n", target.Id, target.Phfm.Velocity)
			}
		}

	case WISH_SPAWN:
		if uint(wish["sponsor"]) == 0 {
			ph.spwan(target)
		}

	case WISH_PING:
		wish["receiver"] = 0
		ph.PostQueue <- wish
	}

}

// moves a cllider without collision check (one of the Phyz hands)
func (ph *Phyz) moveCollider(body *PhyzFm, downward, rightward float64) {
	body.Position[0] += downward
	body.Position[1] += rightward

	//boundary check
	boundI := ph.woldProps["boundary i"]
	boundJ := ph.woldProps["boundary j"]

	switch uint(ph.woldProps["boundary type"]) {

	case BOUNDARY_STATIC:

		if body.Position[0] > boundI {
			body.Position[0] = boundI
		}
		if body.Position[1] > boundJ {
			body.Position[1] = boundJ
		}
		if body.Position[0] < 0 {
			body.Position[0] = 0
		}
		if body.Position[1] < 0 {
			body.Position[1] = 0
		}

	case BOUNDARY_CONTINOUOUS:

		if body.Position[0] > boundI {
			body.Position[0] -= boundI
		}
		if body.Position[1] > boundJ {
			body.Position[1] -= boundJ
		}
		if body.Position[0] < 0 {
			body.Position[0] += boundI
		}
		if body.Position[1] < 0 {
			body.Position[1] += boundJ
		}
	}

}

// checks do a and b have collision (one of the Phyz hands)
func haveCollision(a, b *PhyzFm) bool {

	distance := math.Sqrt(math.Pow(a.Position[0]-b.Position[0], 2) + math.Pow(a.Position[1]-b.Position[1], 2))
	return a.Shape.Properties[0]+a.Shape.Properties[0] > distance
}

// CheckCollision checks global collision
func (ph *Phyz) CheckCollision(target *Object) (iss bool) {

	if !target.Phfm.Present {
		return false
	}

	isInter := ph.interestingColls.Exist(target.Id)
	iss = false

	for iter, other := ph.objBank.GetRoot(); iter != nil; iter, other = ph.objBank.Next(iter) {

		if other.Phfm == nil || other.Id == target.Id || !other.Phfm.Present {
			continue
		}

		if haveCollision(target.Phfm, other.Phfm) {
			iss = true
			if !isInter {
				break
			} else {
				ph.PostQueue <- Wish{
					"type":     float64(WISH_INTERCOLL),
					"receiver": 2,
					"id":       float64(target.Id),
					"othersId": float64(other.Id),
				}
			}
		}
	}
	// fmt.Println("______________________________")
	// fmt.Printf("target: <%v> %v (%+v)  => %v\n", target.Id, target.Phfm.Present, target.Phfm.Position,
	// 	iss)
	return iss
}

func (ph *Phyz) checkDistances() {

	for iter, intl := ph.dataBase.Inteligents.GetRoot(); intl != nil; iter, intl = ph.dataBase.Inteligents.Next(iter) {
		for iter, other := ph.objBank.GetRoot(); other != nil; iter, other = ph.objBank.Next(iter) {
			if other.Id == intl.Id {
				continue
			}

			relPos := cartesianToPolar([2]float64{
				other.Phfm.Position[0] - intl.Phfm.Position[0],
				other.Phfm.Position[1] - intl.Phfm.Position[1]})

			if relPos[0] < intl.Brain.GetSencRadius()[1] {
				w := Wish{
					"type":     WISH_ALART,
					"receiver": 2,
					"id":       float64(intl.Id),
					"othersId": float64(other.Id),
					"radius":   relPos[0],
					"theta":    relPos[1],
				}
				// fmt.Println("Phyz >> alart posted!")
				ph.PostQueue <- w
			}
		}
	}
}

func (ph *Phyz) spinn() {
	for ph.alive {

		<-time.After(ph.timeResultion)
		for iter, mover := ph.moversCache.GetRoot(); mover != nil; iter, mover = ph.moversCache.Next(iter) {

			t := float64(ph.timeResultion/time.Millisecond) / 1000
			dir := polarToCartasian([2]float64{
				mover.Phfm.Velocity[0] * t, mover.Phfm.Velocity[1]})

			ph.moveCollider(mover.Phfm, dir[0], dir[1])
			if ph.CheckCollision(mover) {
				ph.moveCollider(mover.Phfm, -dir[0], -dir[1])
			}

		}
	}
}
