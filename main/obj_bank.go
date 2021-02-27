package main

import (
	"fmt"
	"log"
	"sync"
)

type obj_container struct {
	data *Object
	next *obj_container
}

type ObjList struct {
	root *obj_container
	end  *obj_container
	lock MinionLock
}

func (bn *ObjList) Add(data *Object) {

	if _, founded := bn.Find(data); founded {
		return
	}

	bn.lock.Lock()

	if bn.root == nil {
		bn.root = new(obj_container)
		bn.end = bn.root
	} else {
		bn.end.next = new(obj_container)
		bn.end = bn.end.next
	}
	bn.end.data = data
	bn.end.next = nil

	bn.lock.Unlock()
}

func (bn *ObjList) AddCopy(data *Object) {
	tempPtr := new(Object)
	*tempPtr = *data
	bn.Add(tempPtr)
}

func (bn *ObjList) Get(id uint) (obj *Object, ok bool) {

	bn.lock.MinLock()

	cursor := bn.root
	ok = false
	for cursor != nil {
		if cursor.data.Id == id {
			ok = true
			break
		}
		cursor = cursor.next
	}

	if ok {
		obj = cursor.data
	}
	bn.lock.MinUnlock()
	return
}

func (bn *ObjList) Del(id uint) (ok bool) {

	if bn.root == nil {
		return false
	}

	if bn.root.data.Id == id {
		bn.root = bn.root.next
		return true
	}

	ok = false

	bn.lock.Lock()
	cursor := bn.root
	for cursor.next.next != nil {
		if cursor.next.data.Id == id {
			ok = true
			break
		}
		cursor = cursor.next
	}

	if ok {
		cursor.next = cursor.next.next
	}
	bn.lock.Unlock()
	return ok
}

func (bn *ObjList) Find(target *Object) (id uint, founded bool) {

	bn.lock.MinLock()
	founded = false
	for iter, value := bn.GetRoot(); iter != nil; iter, value = bn.Next(iter) {
		if target == value {
			founded = true
			bn.lock.MinUnlock()
			return
		}
	}
	bn.lock.MinUnlock()
	return
}

func (bn *ObjList) Next(current *obj_container) (address *obj_container, obj *Object) {

	bn.lock.MinLock()
	if current == nil {
		log.Println("ObjBank>> Next was called by nil current")
	}

	if current.next == nil {
		bn.lock.MinUnlock()
		return nil, nil
	}
	address, obj = current.next, current.next.data
	bn.lock.MinUnlock()
	return
}

//GetRoot retruns first elemnts address and context
func (bn *ObjList) GetRoot() (address *obj_container, obj *Object) {

	bn.lock.MinLock()
	if bn.root == nil {
		// log.Println("ObjBank>> GetRoot was called while ObjList was empty")
		bn.lock.MinUnlock()
		return nil, nil
	}
	address, obj = bn.root, bn.root.data
	bn.lock.MinUnlock()
	return
}

func (bn *ObjList) Do(foo func(*Object)) {
}

func (bn *ObjList) PrintMe() {

	bn.lock.MinLock()
	fmt.Printf("\n+-----------------------------+\n")
	for iter, value := bn.GetRoot(); iter != nil; iter, value = bn.Next(iter) {
		fmt.Printf("<%v>  %v ,", value.Id, value.Phfm.Present)
		// fmt.Printf("<%v> ,", value)
	}
	fmt.Println()
	bn.lock.MinUnlock()
}

func (bn *ObjList) Exist(targetId uint) (iss bool) {

	bn.lock.MinLock()
	iss = false
	for iter, value := bn.GetRoot(); iter != nil; iter, value = bn.Next(iter) {
		if targetId == value.Id {
			iss = true
			bn.lock.MinUnlock()
			return
		}
	}
	bn.lock.MinUnlock()
	return
}

type MinionLock struct {
	mainLock     sync.Mutex
	tempLock     sync.Mutex
	minionsCount int
}

func (lk *MinionLock) MinLock() {

	lk.tempLock.Lock()
	if lk.minionsCount == 0 {
		lk.mainLock.Lock()
	}
	lk.minionsCount++
	lk.tempLock.Unlock()
}

func (lk *MinionLock) MinUnlock() {

	lk.tempLock.Lock()
	if lk.minionsCount == 1 {
		lk.mainLock.Unlock()
	}
	lk.minionsCount--
	lk.tempLock.Unlock()
}

func (lk *MinionLock) Lock() {

	lk.mainLock.Lock()
}

func (lk *MinionLock) Unlock() {

	lk.mainLock.Unlock()
}

//Cdb : cneteral data base
type Cdb struct {
	MainList    ObjList
	Movers      ObjList
	Inteligents ObjList
}
