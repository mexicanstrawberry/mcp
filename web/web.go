package web

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
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
		port = "80"
	}

	// load server addr
	var addr string
	if addr = os.Getenv("CF_INSTANCE_ADDR"); len(addr) == 0 {

		// CF_INSTANCE_ADDR is empty, set it to localhost
		addr = ":" + port

	} else {
		// CF_INSTANCE_ADDR is set
		addr = ":" + port
	}

	r := mux.NewRouter()

	r.HandleFunc("/", index)

	r.HandleFunc("/api/{version}/service/{id:.+}/start_recipe", ServiceStartRecipe)

	http.Handle("/", r)

	go func() {
		// start server on port 5000 for api request
		apiAddr := ":5000"
		clog.Infoln("Server Listening", "address", apiAddr)
		http.ListenAndServe(":80", nil)
	}()

	go RunCouchProxy()

	clog.Infoln("Server Listening", "address", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		clog.Fatal(err)
	}

}
