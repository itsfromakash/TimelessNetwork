package main

import "C"

import (
	"TimelessNetwork/P2P"
)

//export StartP2PEngine
func StartP2PEngine() {
	P2P.StartP2PEngine()
}

func main() {}