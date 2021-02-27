package main

import (
	"reflect"
)

func newDumC(pos [2]float64) (ng *Object) {
	ng = new(Object)
	ng.Phfm = new(PhyzFm)
	ng.Phfm.Position = pos
	ng.Phfm.Shape = &GeoShape{
		Typee:      "circle",
		Properties: []float64{1},
	}

	ng.PostQueue = make(chan Wish, 10)
	return
}

func isNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}

// func newLooser(pos [2]float64) (ng *Object) {
// 	ng = newDumC(pos)
// 	ng.Brain =

// }
