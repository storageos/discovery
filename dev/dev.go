package main

import (
	"log"

	"github.com/rsc/devweb/slave"
	_ "github.com/storageos/discovery/http"
)

func main() {
	log.SetFlags(0)
	slave.Main()
}
