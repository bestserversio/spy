package utils

import "fmt"

func DebugMsg(reqLvl int, curLvl int, msg string, args ...interface{}) {
	if curLvl >= reqLvl {
		fmt.Printf(msg+"\n", args)
	}
}
