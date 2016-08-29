package main

import (
	"fmt"
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

	addr := "0.0.0.0"

	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = "8080"
	}

	listenAddr := fmt.Sprintf("%s:%s", addr, port)
	clog.Infoln("Server Listening", "address", listenAddr, "protocol")

	http.HandleFunc("/", index)
	err = http.ListenAndServe(listenAddr, nil)
	if err != nil {
		clog.Fatal(err)
	}

}
