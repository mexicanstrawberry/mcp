package web

import (
	"os"
	"time"

	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/mexicanstrawberry/go-cfenv"
	clog "github.com/morriswinkler/cloudglog"
)

var (
	couchDbURL     string
	couchProxyPort = ":5984"
)

func Init() {

	// this function is fird if VCAP_SERVICES was not found and refires every 10 Seconds
	var vcapServicesNotFound func() // declare first to access var recursively
	vcapServicesNotFound = func() {

		clog.Errorln("[web][cochdbProxy] ###############################################")
		clog.Errorln("[web][cochdbProxy] CRITICAL ERROR: could not read VCAP_SERVICES   ")
		clog.Errorln("[web][cochdbProxy] ###############################################")
		for _, k := range os.Environ() {
			clog.Errorln(k)
		}

		time.AfterFunc(time.Second*10, vcapServicesNotFound) // recursion

	}

	appEnv, err := cfenv.Current()
	if err != nil {
		// this is critical VCAP_SERVICES was not found
		vcapServicesNotFound()
	} else {

		service, err := appEnv.Services.WithName("Cloudant NoSQL DB-hg")
		if err != nil {
			clog.Errorln("[web][cochdbProxy] ", err)

		} else {

			couchDbURL = service.Credentials["url"].(string)
		}
	}

}

//  our RerverseProxy object
type CouchProx struct {
	target *url.URL               // target url of reverse proxy
	proxy  *httputil.ReverseProxy // instance of Go ReverseProxy
	mux    *http.ServeMux         // mux http handler
}

// small factory
func newCouchProxy(target string) *CouchProx {
	url, _ := url.Parse(target)
	// you should handle error on parsing

	couchProxy := &CouchProx{
		target: url,
		proxy:  httputil.NewSingleHostReverseProxy(url),
		mux:    http.NewServeMux(),
	}

	return couchProxy
}

func (p *CouchProx) handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-GoProxy", "GoProxy")
	// call to magic method from ReverseProxy object
	p.proxy.ServeHTTP(w, r)
}

func RunCouchProxy() {

	Init()

	clog.Infoln("Starting CouchDbProxy address:", couchProxyPort)

	// proxy
	proxy := newCouchProxy(couchDbURL)

	// handler
	proxy.mux.HandleFunc("/", proxy.handle)

	// server
	http.ListenAndServe(couchProxyPort, proxy.mux)
}
