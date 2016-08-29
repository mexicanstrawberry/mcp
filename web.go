package main

import (
	"net/http"
	"os"

	clog "github.com/morriswinkler/cloudglog"
)

func index(w http.ResponseWriter, r *http.Request) {

	msg := "Mexicanstrawberry is online"

	w.Write([]byte(msg))
}

func StartHttpServer() {

	var err error

	/*
	* Load server address and port
	 */
	// load port
	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = "8080"
	}

	// load server addr
	var addr string
	if addr = os.Getenv("CF_INSTANCE_ADDR"); len(addr) == 0 {

		// CF_INSTANCE_ADDR is empty, set it to localhost
		addr = "127.0.0.1:" + port

	} else {
		// CF_INSTANCE_ADDR is set
		addr = ":" + port
	}

	clog.Infoln("Server Listening", "address", addr, "protocol")

	http.HandleFunc("/", index)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		clog.Fatal(err)
	}

}
