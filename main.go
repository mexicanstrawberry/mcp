package main

import (
	"flag"

	"github.com/mexicanstrawberry/mcp/gatekeeper"
	"github.com/mexicanstrawberry/mcp/recipe"
	"github.com/mexicanstrawberry/mcp/sensor"
	"github.com/mexicanstrawberry/mcp/web"
	clog "github.com/morriswinkler/cloudglog"
)

var (
	debug             int    // used to set debug mode 1 2 or 3
	dumpCouchDbServer string // dump couchdb server
	localhost         bool   // run on localhost, loads some prepared ssl certs
)

func init() {

	flag.IntVar(&debug, "debug", 0, "debug level 1,2 ,3")
	flag.StringVar(&dumpCouchDbServer, "dumpDB", "", "dump couchdb server <[user:pass@]url[:port]>")
	flag.BoolVar(&localhost, "localhost", false, "load prepared dummy ssl certs for localhosts")

	flag.Parse()
}

func main() {

	// handle debug values
	switch debug {
	case 3:
		gatekeeper.DEBUG = true
		clog.LogLevel = 3
	case 2:
		clog.LogLevel = 2
	case 1:
		clog.LogLevel = 1
	}

	if dumpCouchDbServer != "" {
		clog.Infoln("dumping couchdb database: ", dumpCouchDbServer)
		err := web.DumpCouchDB(dumpCouchDbServer)
		if err != nil {
			clog.Errorln("[backupCouchDB]", err)
		}
	}

	// TODO: implement debug clog.LogLevel and gatekeeper.Debug

	clog.Infoln("init mexicanstrawberry")

	myRecipe, err := recipe.NewRecipe("SimpleRecipe", true) // true for dummy recipe
	if err != nil {
		clog.Errorln(err)
	}

	insideHumidity, err := sensor.NewSensor("InsideHumidity", myRecipe)
	if err != nil {
		clog.Errorln(err)
	}

	insideTemperature, err := sensor.NewSensor("InsideTemperature", myRecipe)
	if err != nil {
		clog.Errorln(err)
	}

	go gatekeeper.Run()
	go insideHumidity.Run()
	go insideTemperature.Run()

	web.StartHttpServer()

}
