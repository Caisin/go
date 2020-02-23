package main

import (
	"log"
	"testing"
)

func TestIota(t *testing.T) {
	type Flags uint
	const (
		FlagUp           Flags = 1 << iota // is up
		FlagBroadcast                      //supports broadcast access capability
		FlagLoopback                       // is a loopback interface
		FlagPointToPoint                   // belongs to a point-to-point link
		FlagMulticast                      // supports multicast access capability
	)
	log.Print(FlagUp, FlagBroadcast, FlagLoopback, FlagPointToPoint, FlagMulticast)
	r := [...]int{99: -1}
	log.Print(r)
}
