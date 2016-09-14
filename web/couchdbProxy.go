package web

import (
	"os"
	"time"

	"net/http"
	"net/http/httputil"
	"net/url"

	"io/ioutil"
	"strings"

	"encoding/json"

	"path/filepath"

	"github.com/mexicanstrawberry/go-cfenv"
	clog "github.com/morriswinkler/cloudglog"
)

var (
	couchDbURL     string
	couchProxyPort = ":5984"
)

func init() {

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

	clog.Infoln("Starting CouchDbProxy address:", couchProxyPort)

	// proxy
	proxy := newCouchProxy(couchDbURL)

	// handler
	proxy.mux.HandleFunc("/", proxy.handle)

	// server
	ListenAndServeTLS(couchProxyPort, proxy.mux)
}

// DumpCouchDb takes a url and dumps all databases behind it
// expect a few blacklisted ones like _users and _replicator
func DumpCouchDB(baseurl string) error {

	// create http client for reuse
	client := http.Client{
		Timeout: time.Second * 10,
	}

	// get all databases
	resp, err := client.Get(concatURL(baseurl, "_all_dbs"))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	clog.V(1).Infoln("[web][DumpCouchDB]", string(body))

	var databases []string

	err = json.Unmarshal(body, &databases)
	if err != nil {
		return err
	}

	clog.Infoln("[web][DumpCouchDB] backing up:", databases)

	// black listed db names, do not dump
	var dbBlacklist = []string{"_users", "_replicator", "test_suite_db2"}

	for _, db := range databases {

		// A "continue" statement begins the next iteration of the innermost "for" loop
		// this is a ugly hack to overcome this behaviour
		if func() bool {
			for _, val := range dbBlacklist {
				if db == val {
					return true
				}
			}
			return false
		}() {
			continue
		}

		respDB, err := client.Get(concatURL(baseurl, db, "_all_docs?include_docs=true&attachments=true"))

		body, err := ioutil.ReadAll(respDB.Body)
		if err != nil {
			return err
		}

		clog.V(2).Infoln(string(body))

		err = dumpDBtoFile(db, body)
		if err != nil {
			return err
		}

		respDB.Body.Close()
	}

	return nil

}

func concatURL(baseurl string, endpoints ...string) string {

	var result string

	if strings.HasSuffix(baseurl, "/") {
		endpoints = append([]string{baseurl[:len(baseurl)-1]}, endpoints...)
	} else {
		endpoints = append([]string{baseurl}, endpoints...)
	}

	result = strings.Join(endpoints, "/")

	clog.V(1).Infoln("[web][DumpCouchDB]", result)

	return result
}

// save database json string to a file
func dumpDBtoFile(db string, dbDump []byte) error {

	wdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	dir := strings.Join([]string{wdir, "couchDb_dumps"}, "/")

	clog.V(2).Infoln("[web][DumpCouchDB] dumping to directory", wdir)

	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	fileName := strings.Join([]string{dir, db}, "/")
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// see https://github.com/danielebailo/couchdb-dump/blob/master/couchdb-backup.sh
	// CouchDB has a tendancy to output Windows carridge returns in it's output -
	// This messes up us trying to sed things at the end of lines!

	dbDumpString := strings.NewReplacer("\r\n", "\r\n", "\n", "\r\n").Replace(string(dbDump))

	_, err = file.Write([]byte(dbDumpString))
	if err != nil {
		return err
	}

	clog.V(2).Infoln("[web][DumpCouchDB] dumping to file", fileName)

	return nil
}

//TODO: to be impemented
func uploadDBfromDump(db string, dbDump []byte) error {

	// seperate design files from dbDump
	// grep '^{"id":"_design' ${file_name} > ${design_file_name}

	// cut design files from dbDump
	// sed ${sed_edit_in_place} '/^{"_id":"_design/d' ${file_name} && rm -f ${file_name}.sedtmp

	// remove the trailing comma from dbDump
	// sed ${sed_edit_in_place} "${line}s/,$//" ${file_name} && rm -f ${file_name}.sedtmp

	// Split the ID out for use as the import URL path
	// URLPATH=$(echo $line | awk -F'"' '{print$4}')

	// Scrap the ID and Rev from the main data, as well as any trailing ','
	// echo "${line}" | sed -${sed_regexp_option}e "s@^\{\"_id\":\"${URLPATH}\",\"_rev\":\"[0-9]*-[0-9a-zA-Z_\-]*\",@\{@" | sed -e 's/,$//' > ${design_file_name}.${designcount}

	// Search for \\n and remove

	// uplaod design files
	// curl -T ${design_file_name}.${designcount} -X PUT "${url}/${db_name}/${URLPATH}" -H 'Content-Type: application/json' -o ${design_file_name}.out.${designcount}

	// upload dbDump
	// curl -T $file_name -X POST "$url/$db_name/_bulk_docs" -H 'Content-Type: application/json' -o tmp.out

	return nil
}
