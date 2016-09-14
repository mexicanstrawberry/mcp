package web

import (
	"net/http"
	"os"

	"strings"

	"github.com/gorilla/mux"
	clog "github.com/morriswinkler/cloudglog"
)

func index(w http.ResponseWriter, r *http.Request) {

	msg := "Mexicanstrawberry is online"

	w.Write([]byte(msg))
}

var port, tlsPort string

func StartHttpServer() {

	var err error

	/*
	* Load server address and port
	 */
	// load port
	if port = os.Getenv("PORT"); len(port) == 0 {
		// standart ports
		port = "80"
		tlsPort = "443"
	} else {
		// custom port ( http 8080 -> https 8084 )
		tlsPortB := []byte(port)                      // convert to byte slice for char replacemnet
		copy(tlsPortB[len(tlsPortB)-1:], []byte("4")) // replace last digit with 4
		tlsPort = string(tlsPortB)                    // convert back to string
	}

	// load server addr
	var addr, tlsAddr string
	if addr = os.Getenv("CF_INSTANCE_ADDR"); len(addr) == 0 {

		// CF_INSTANCE_ADDR is empty, set it to localhost
		addr = ":" + port
		tlsAddr = ":" + tlsPort

	} else {
		// CF_INSTANCE_ADDR is set
		addr = ":" + port
	}

	/*
	 * define route handlers
	 */

	r := mux.NewRouter()

	r.HandleFunc("/", index)

	r.HandleFunc("/api/{version}/service/{id:.+}/start_recipe", ServiceStartRecipe)

	http.Handle("/", r)

	/*
	 * start server
	 */

	// start couchDb proxy
	go RunCouchProxy()

	go func() {
		// start server on port 5000 for api request
		apiAddr := ":5000"

		clog.Infoln("API Server Listening", "address", apiAddr)

		err = ListenAndServeTLS(apiAddr, nil)
		if err != nil {
			clog.Fatal(err)
		}
	}()

	go func() {
		// start server on port 80 and redirect all request to the TLS port
		rMux := http.NewServeMux()
		rMux.HandleFunc("/", RedirectHTTP)

		clog.Infoln("Starting redirect from ", addr, "to", tlsAddr)

		err = http.ListenAndServe(addr, rMux)
		if err != nil {
			clog.Fatal(err)
		}
	}()

	clog.Infoln("Server Listening", tlsAddr)

	// start TLS web server
	err = ListenAndServeTLS(tlsAddr, nil)
	if err != nil {
		clog.Fatal(err)
	}

}

func RedirectHTTP(w http.ResponseWriter, r *http.Request) {
	if r.TLS != nil || r.Host == "" {
		http.Error(w, "not found", 404)
	}

	u := r.URL

	// add tlsPort :XXXX to u.Host
	// if tlsPort is not equal 443
	var host = r.Host
	if tlsPort != "443" {
		// no standart TLS port
		host = strings.TrimSuffix(r.Host, port) // remove trailing port
		host = strings.Join([]string{host, tlsPort}, "")
	}

	u.Host = host
	u.Scheme = "https"
	http.Redirect(w, r, u.String(), 302)
}
