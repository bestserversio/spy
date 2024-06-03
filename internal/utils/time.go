package utils

import (
	"math/rand"
	"time"
)

func RandomWait(min int, max int) {
	wait := rand.Intn(max-min+1) + min

	time.Sleep(time.Second * time.Duration(wait))
}
