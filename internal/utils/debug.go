package utils

import "fmt"

func DebugMsg(reqLvl int, curLvl int, msg string, args ...interface{}) {
	if curLvl >= reqLvl {
		if len(args) > 0 {
			fmt.Printf(msg+"\n", args...)
		} else {
			fmt.Printf(msg + "\n")
		}
	}
}
