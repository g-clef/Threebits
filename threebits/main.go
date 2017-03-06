package main

import (
	"github.com/accessviolationsec/Threebits"
	"math/rand"
	"time"
)

func init(){
	rand.Seed(time.Now().UTC().UnixNano())
}

func main(){
	Threebits.Run()
}
