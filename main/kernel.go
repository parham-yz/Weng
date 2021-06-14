package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type Kernel struct {
	HeaderFm
	dataBase       *Cdb
	serversHeaders []*HeaderFm
	serveQueue     chan Wish
	idGenerator    uint
	peng           *Phyz
	crier          *Crier
	alive          bool
}

func (kr *Kernel) Init(dataBase *Cdb) {
	if kr.peng == nil {
		kr.peng = new(Phyz).Init(dataBase, map[string]float64{
			"boundary type": BOUNDARY_STATIC,
			"boundary i":    10,
			"boundary j":    10,
		})
	}

	if kr.crier == nil {
		kr.crier = new(Crier).Init(dataBase)
	}

	kr.dataBase = dataBase
	kr.PostQueue = make(chan Wish, 100)
	kr.serveQueue = make(chan Wish, 100)
	kr.Id = 0
	kr.idGenerator = 10
	kr.serversHeaders = []*HeaderFm{&kr.HeaderFm, &kr.peng.HeaderFm, &kr.crier.HeaderFm}

	os.Remove("log.txt")
	logFile, _ := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY, 0666)
	log.SetOutput(logFile)

}

func (kr *Kernel) DropIn(obj *Object, tag string) *HeaderFm {
	if !kr.alive {
		panic("Krenel is dead")
	}

	objc := new(Object)
	*objc = *obj

	objc.Id = kr.idGenerator
	kr.idGenerator++
	objc.Tag = tag

	kr.dataBase.MainList.Add(objc)

	kr.peng.PostQueue <- Wish{
		"receiver": 1,
		"id":       float64(objc.Id),
		"sponsor":  float64(kr.Id),
		"type":     WISH_SPAWN,
	}

	if !isNil(obj.Brain) {
		kr.dataBase.Inteligents.Add(objc)
	}

	if obj.Phfm.Velocity[0] != 0 {
		kr.dataBase.Movers.Add(objc)
	}

	return &objc.HeaderFm
}

func (kr *Kernel) PullOut(id uint) {
	if !kr.alive {
		panic("Krenel is dead")
	}

	if kr.dataBase.Inteligents.Exist(id) {
		kr.dataBase.Inteligents.Del(id)
	}

	if kr.dataBase.Movers.Exist(id) {
		kr.dataBase.Movers.Del(id)
	}

	kr.dataBase.MainList.Del(id)
}

func (kr *Kernel) RouteWishes() {

	// log.Println("kernel >> routing wishes")
	// routing servers wishes
	for _, server := range kr.serversHeaders {
		for i := 0; i < 100; i++ {
			select {
			case wish := <-server.PostQueue:
				kr.PostToServers(wish)
			default:
				i += 25
			}
		}
	}

	//Objs wishes :
	for iter, target := kr.dataBase.MainList.GetRoot(); target != nil; iter, target = kr.dataBase.MainList.Next(iter) {

		for i := 0; i < 10; i++ {
			select {
			case wish := <-target.PostQueue:
				kr.PostToServers(wish)
			default:
				i += 10
			}
		}
	}
}

func (kr *Kernel) PostToServers(wish Wish) {
	if kr.peng == nil || kr.crier == nil {
		panic("missing server!")
	}

	// log.Printf("Kernel >> post<%v,%v> to %v\n", wish["id"], wish["type"], wish["receiver"])
	switch uint(wish["receiver"]) {
	case 0:
		kr.serveQueue <- wish
	case 1:
		kr.peng.serveQueue <- wish
	case 2:
		kr.crier.serveQueue <- wish
	}
}

func (kr *Kernel) serveWish(wish Wish) {
	switch wish["type"] {
	case WISH_PING:
		log.Print("Ping! : ")
		for key, value := range wish {
			log.Printf(" %v: %v ", key, value)
		}
	}

}

func (kr *Kernel) Run() {
	for {

		select {
		case wish, chanOpen := <-kr.serveQueue:
			if !chanOpen {
				log.Println("killing Kernel")
				return
			}
			kr.serveWish(wish)
		case <-time.After(time.Second / 10):
			kr.RouteWishes()

		}
	}

}

func (kr *Kernel) StartIt(withPortal bool) {
	kr.alive = true
	go kr.Run()
	go kr.peng.Run()
	go kr.crier.Run()

	if withPortal {
		go OpenPortall(&kr.dataBase.MainList, &kr.alive, 16, time.Second/4)
	}
	log.Println("kernel started")
}

func (kr *Kernel) Shutdown() {
	close(kr.serveQueue)
	close(kr.peng.serveQueue)
	close(kr.crier.serveQueue)

	kr.alive = false
}

//OpenPortall lets you see your wold form a baounded window
func OpenPortall(bn *ObjList, alive *bool, windowSize float64, frameTime time.Duration) {
	var vmap [40][90]uint8

	charSize := [2]float64{float64(windowSize / 40), float64(windowSize / 90)}

	for *alive {
		<-time.After(frameTime)
		for iter, target := bn.GetRoot(); target != nil; iter, target = bn.Next(iter) {
			if target.Phfm.Position[0] < windowSize && target.Phfm.Position[1] < windowSize {
				index := [2]int{int(target.Phfm.Position[0] / charSize[0]), int(target.Phfm.Position[1] / charSize[1])}
				vmap[index[0]][index[1]] = 'O'
			}
		}
		printCharMap(vmap)

		for i, row := range vmap {
			for j := range row {
				vmap[i][j] = ' '
			}
		}
	}

}

func printCharMap(mp [40][90]uint8) {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	for _, row := range mp {
		for _, pos := range row {
			fmt.Printf("%c", pos)
		}
		fmt.Println("|")
	}

	for range mp[0] {
		fmt.Print("-")
	}
	fmt.Println()

}
