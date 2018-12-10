package main

import (
	"github.com/wlibo666/common-lib/webutils"
)

func init() {
	webutils.AddController(samplerGetController)
	webutils.AddController(samplerSetController)
}

func main() {
	webutils.ServerRun(":10100")
}
