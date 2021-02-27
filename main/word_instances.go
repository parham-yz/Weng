package main

import (
	"math/rand"
)

const (
	WISH_MOVE = iota
	WISH_SPAWN
	WISH_KILL
	WISH_ALART
	WISH_PING
	WISH_INTERCOLL
)

const (
	SIG_ALART = iota
	SIG_COLL
)

type Wish map[string]float64

type HeaderFm struct {
	Id        uint
	Tag       string
	PostQueue chan Wish
}

type PhyzFm struct {
	Present  bool
	Position [2]float64
	Shape    *GeoShape
	Velocity [2]float64
}

type Intel interface {
	Decide() [2]float64
	PushEnvSig(*EnvSig)
	GetSencRadius() [2]float64
}

type EnvSig struct {
	Type        int
	RelPolarPos [2]float64
	Reliability float64
	OthersTag   string
}

type Object struct {
	HeaderFm
	Phfm  *PhyzFm
	Brain Intel
}

func (hd *HeaderFm) ping(receiverID uint) {
	hd.PostQueue <- Wish{
		"receiver": float64(receiverID),
		"id":       float64(hd.Id),
		"type":     WISH_PING,
	}
}

func (obj *Object) MoveMe(radius, theta float64) {

	w := Wish{
		"receiver": 1,
		"id":       float64(obj.Id),
		"type":     WISH_MOVE,
		"radius":   radius,
		"theta":    theta,
	}
	obj.PostQueue <- w
}

func (obj *Object) RandomWalk(scale float64) {

	w := Wish{
		"receiver": 1,
		"id":       float64(obj.Id),
		"type":     WISH_MOVE,
		"dw":       (1 - 2*rand.Float64()) * scale,
		"rw":       (1 - 2*rand.Float64()) * scale,
	}
	obj.PostQueue <- w
}
