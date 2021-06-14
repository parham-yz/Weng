package main

import (
	"fmt"
	"time"
)

var kernel Kernel
var peng Phyz
var db Cdb

func main() {

	kernel.peng = peng.Init(&db, map[string]float64{
		"boundary type": BOUNDARY_CONTINOUOUS,
		"boundary i":    16,
		"boundary j":    16,
	})
	kernel.Init(&db)
	fmt.Println("starting")
	kernel.StartIt(true)

	jimi := newDumC([2]float64{0, 0})
	walls := make([]*Object, 0)
	for i := 0.0; i < 15; i += 2 {
		walls = append(walls, newDumC([2]float64{7, i}))

	}

	for _, wall := range walls {
		kernel.DropIn(wall, "")
	}

	jimi.Brain = new(gongishkBrain)
	jimi.HeaderFm = *kernel.DropIn(jimi, "")
	jimi.MoveMe(1, PI)

	for i := 0; i < 50; i++ {
		// fmt.Println(jimi.Phfm.Position)
		time.Sleep(time.Second / 2)
	}

	time.Sleep(time.Second * 2)
	kernel.Shutdown()
	println("end")

}
