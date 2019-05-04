package main

import (
	"github.com/calebcase/core"
	"github.com/inconshreveable/log15"
)

func main() {
	core.Log.SetHandler(log15.StderrHandler)
	core.DumpAll()
}
