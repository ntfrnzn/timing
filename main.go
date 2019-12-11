package main

import (
	"time"

	"github.com/ntfrnzn/timing/pkg/util"
)

func test1() {
	if util.FuncTimer != nil {
		timer := util.FuncTimer.Instrument()
		defer timer()
	}
	time.Sleep(2 * time.Second)
}

func test2() {
	if util.FuncTimer != nil {
		timer := util.FuncTimer.Instrument()
		defer timer()
	}
	time.Sleep(4 * time.Second)
}

func main() {

	if util.FuncTimer != nil {
		go util.FuncTimer.Receive()
	}

	go test1()
	go test1()
	test2()

	if util.FuncTimer != nil {
		util.FuncTimer.Terminate()
	}
}
